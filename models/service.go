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
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Upstreams []string `json:"upstreams"`

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
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			// Maybe we should prefix the service buckets instead of excluding other known ones.
			if string(name) == "Services" {
				return nil
			}

			service, err := getServiceByName(name)
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
	s, err := maybeGetServiceByName([]byte(service.Name))
	if err != nil {
		return nil, err
	}

	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	if s == nil {
		// Create a new service
		s = &Service{
			Name:      service.Name,
			Path:      service.Path,
			Upstreams: service.Upstreams,

			Registered: time.Now(),
		}

		err := db.Update(func(tx *bolt.Tx) error {
			serviceBucket, err := tx.CreateBucket([]byte(s.Name))
			if err != nil {
				log.Error(err)
				return err
			}

			// Write the path
			if err := serviceBucket.Put([]byte("path"), []byte(s.Path)); err != nil {
				log.Error(err)
				return err
			}

			// Write the upstreams
			for _, upstream := range service.Upstreams {
				log.Debugf("Creating upstream %q for service %q", upstream, s.Name)
				if err := serviceBucket.Put([]byte(fmt.Sprintf("upstream:%s:url", upstream)), []byte(upstream)); err != nil {
					log.Error(err)
					return err
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Update the existing service
		s.Name = service.Name

		err := db.Update(func(tx *bolt.Tx) error {
			serviceBucket := tx.Bucket([]byte(s.Name))

			// Update the path
			if err := serviceBucket.Put([]byte("path"), []byte(s.Path)); err != nil {
				log.Error(err)
				return err
			}

			// Don't create multiple upstreams with the same url, so check if it exists
			toCreate := utils.DiffArrays(service.Upstreams, s.Upstreams)
			for _, upstream := range toCreate {
				log.Debugf("Adding upstream %q to service %q", upstream, service.Name)
				if err := serviceBucket.Put([]byte(fmt.Sprintf("upstream:%s:url", upstream)), []byte(upstream)); err != nil {
					log.Error(err)
					return err
				}

				s.Upstreams = append(s.Upstreams, upstream)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func maybeGetServiceByName(name []byte) (*Service, error) {
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	service := Service{
		Name: string(name),
	}

	ok := false
	err = db.View(func(tx *bolt.Tx) error {
		serviceBucket := tx.Bucket(name)
		if serviceBucket == nil {
			return nil
		}

		ok = true
		service.Path = string(serviceBucket.Get([]byte("path")))

		service.Upstreams = make([]string, 0, 0)
		err := serviceBucket.ForEach(func(k, v []byte) error {
			if strings.HasPrefix(string(k), "upstream:") && strings.HasSuffix(string(k), ":url") {
				service.Upstreams = append(service.Upstreams, string(v))
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
