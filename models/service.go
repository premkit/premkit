package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/premkit/premkit/log"
	"github.com/premkit/premkit/persistence"
	"github.com/premkit/premkit/utils"

	"github.com/boltdb/bolt"
)

// Service represents a single registered service with this reverse proxy.
type Service struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Upstreams []*Upstream `json:"upstreams"`

	Registered time.Time `json:"registered"`
}

// ListServices returns a list of all available, known services.
// TODO this should cache and not always hit the disk.
func ListServices() ([]*Service, error) {
	services := make([]*Service, 0, 0)

	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(bucketName []byte, b *bolt.Bucket) error {
			if !strings.HasPrefix(string(bucketName), "service:") {
				return nil
			}

			s := strings.Split(string(bucketName), ":")
			if len(s) < 2 {
				err = fmt.Errorf("Unexpected bucket name: %q", bucketName)
				log.Error(err)
				return nil
			}

			serviceName := strings.Join(s[1:], ":")

			service, err := getServiceByName([]byte(serviceName))
			if err != nil {
				return err
			}

			services = append(services, service)

			return nil
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return services, nil
}

// CreateService will create a new (or update an existing) service.  If the service already
// exists, this call will update it with the new name, and append it's own upstream.
// This could be problematic if two different services register with the same path.  The router
// would send traffic randomly to each.
func CreateService(service *Service) (*Service, error) {
	log.Debugf("Creating service %q (path: %q)", service.Name, service.Path)

	if err := validateService(service); err != nil {
		return nil, err
	}

	// Clean the service a little
	service.Path = strings.TrimPrefix(service.Path, "/")

	// If the service already exists, we just want to update it with a new upstream
	current, err := maybeGetServiceByName([]byte(service.Name))
	if err != nil {
		return nil, err
	}

	if current == nil {
		return createNewService(service)
	}

	return updateService(current, service)
}

func updateService(current *Service, service *Service) (*Service, error) {
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		serviceBucket := tx.Bucket([]byte(fmt.Sprintf("service:%s", service.Name)))

		// Update the path
		if err := serviceBucket.Put([]byte("path"), []byte(service.Path)); err != nil {
			log.Error(err)
			return err
		}

		// TODO update the registration date

		// On update, we merge these upstreams into the current upstreams
		upstreamURLs := make([]string, 0, 0)
		for _, u := range current.Upstreams {
			upstreamURLs = append(upstreamURLs, u.URL)
		}
		for _, u := range service.Upstreams {
			upstreamURLs = append(upstreamURLs, u.URL)
		}
		upstreamURLs = utils.RemoveDuplicates(upstreamURLs)

		combinedUpstreams := make([]*Upstream, 0, 0)
		for _, url := range upstreamURLs {
			// Check for this upstream in the new service first
			added := false
			for _, upstream := range service.Upstreams {
				if upstream.URL == url {
					combinedUpstreams = append(combinedUpstreams, upstream)
					added = true
					break
				}
			}

			if !added {
				for _, upstream := range current.Upstreams {
					if upstream.URL == url {
						combinedUpstreams = append(combinedUpstreams, upstream)
						added = true
						break
					}
				}
			}
		}

		service.Upstreams = combinedUpstreams

		for _, upstream := range combinedUpstreams {
			if err := SaveUpstream(upstream, tx); err != nil {
				return err
			}

			// And save the reference
			if err := serviceBucket.Put([]byte(fmt.Sprintf("upstream:%s", upstream.URL)), []byte(upstream.URL)); err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return service, nil
}

func createNewService(service *Service) (*Service, error) {
	// Create a new service
	service.Registered = time.Now()

	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		serviceBucket, err := tx.CreateBucket([]byte(fmt.Sprintf("service:%s", service.Name)))
		if err != nil {
			log.Error(err)
			return err
		}

		// Write the path
		if err := serviceBucket.Put([]byte("path"), []byte(service.Path)); err != nil {
			log.Error(err)
			return err
		}

		// TODO write the registration date

		// Write the upstreams
		for _, upstream := range service.Upstreams {
			// upstream are stored in the service bucket, but these are
			// just references to the upstream buckets themselves.  the details
			// of an upstream must be read from the upstream bucket.
			log.Debugf("Saving upstream with URL %q", upstream.URL)
			if err := SaveUpstream(upstream, tx); err != nil {
				return err
			}

			// And save the reference
			if err := serviceBucket.Put([]byte(fmt.Sprintf("upstream:%s", upstream.URL)), []byte(upstream.URL)); err != nil {
				log.Error(err)
				return err
			}

			log.Debugf("Saved upstream with URL %q", upstream.URL)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return service, nil
}

func maybeGetServiceByName(name []byte) (*Service, error) {
	log.Debugf("Attempting to load a service named %q", name)
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	service := Service{
		Name: string(name),
	}

	ok := false
	err = db.View(func(tx *bolt.Tx) error {
		serviceBucket := tx.Bucket([]byte(fmt.Sprintf("service:%s", name)))
		if serviceBucket == nil {
			return nil
		}

		ok = true
		service.Path = string(serviceBucket.Get([]byte("path")))

		service.Upstreams = make([]*Upstream, 0, 0)
		err := serviceBucket.ForEach(func(k, v []byte) error {
			log.Debugf("key: %q", k)
			if strings.HasPrefix(string(k), "upstream:") {
				upstream, err := maybeGetUpstreamByURL(v)
				if err != nil {
					return err
				}

				if upstream != nil {
					service.Upstreams = append(service.Upstreams, upstream)
				}
			}
			if err != nil {
				log.Error(err)
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	log.Debugf("maybeGetService found a service for name %q", string(name))
	return &service, nil
}

func getServiceByName(name []byte) (*Service, error) {
	service, err := maybeGetServiceByName(name)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, errors.New("Service not found")
	}

	return service, nil
}

func validateService(service *Service) error {
	// TODO this.
	return nil
}
