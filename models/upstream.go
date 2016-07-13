package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/premkit/premkit/log"
	"github.com/premkit/premkit/persistence"

	"github.com/boltdb/bolt"
)

// Upstream represents a single upstream that will be added to a service.
type Upstream struct {
	URL string `json:"url"`

	IncludeServicePath bool `json:"include_service_path"`
	IgnoreInsecure     bool `json:"ignore_insecure"`
}

// SaveUpstream will persist an upstream to the database. This will check the
// database and update the upstream, if it already exists. Upstreams are unique by URL, so
// if the upstream was added to a different service, saving a change will update all services
// that shared this backend.
func SaveUpstream(upstream *Upstream, tx *bolt.Tx) error {
	log.Debugf("Creating or updating upstream %q", upstream.URL)

	existing, err := maybeGetUpstreamByURL([]byte(upstream.URL))
	if err != nil {
		return err
	}

	if existing != nil {
		// Update the current upstream
		upstreamBucket := tx.Bucket([]byte(fmt.Sprintf("upstream:%s", upstream.URL)))
		if upstreamBucket == nil {
			err := errors.New("Bucket not found")
			log.Error(err)
			return err
		}

		// Write the fields
		if err := upstreamBucket.Put([]byte("url"), []byte(upstream.URL)); err != nil {
			log.Error(err)
			return err
		}

		if err := upstreamBucket.Put([]byte("include.service.path"), []byte(strconv.FormatBool(upstream.IncludeServicePath))); err != nil {
			log.Error(err)
			return err
		}

		if err := upstreamBucket.Put([]byte("ignore.insecure"), []byte(strconv.FormatBool(upstream.IgnoreInsecure))); err != nil {
			log.Error(err)
			return err
		}

		return nil
	}

	// Register a new upstream
	upstreamBucket, err := tx.CreateBucket([]byte(fmt.Sprintf("upstream:%s", upstream.URL)))
	if err != nil {
		log.Error(err)
		return err
	}

	// Write the fields
	if err := upstreamBucket.Put([]byte("url"), []byte(upstream.URL)); err != nil {
		log.Error(err)
		return err
	}

	if err := upstreamBucket.Put([]byte("include.service.path"), []byte(strconv.FormatBool(upstream.IncludeServicePath))); err != nil {
		log.Error(err)
		return err
	}

	if err := upstreamBucket.Put([]byte("ignore.insecure"), []byte(strconv.FormatBool(upstream.IgnoreInsecure))); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func maybeGetUpstreamByURL(url []byte) (*Upstream, error) {
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	upstream := Upstream{
		URL: string(url),
	}

	ok := false
	err = db.View(func(tx *bolt.Tx) error {
		upstreamBucket := tx.Bucket([]byte(fmt.Sprintf("upstream:%s", url)))
		if upstreamBucket == nil {
			return nil
		}

		ok = true

		b, err := strconv.ParseBool(string(upstreamBucket.Get([]byte("include.service.path"))))
		if err != nil {
			log.Error(err)
			return err
		}
		upstream.IncludeServicePath = b

		b, err = strconv.ParseBool(string(upstreamBucket.Get([]byte("ignore.insecure"))))
		if err != nil {
			log.Error(err)
			return err
		}
		upstream.IgnoreInsecure = b

		return nil
	})

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &upstream, nil
}
