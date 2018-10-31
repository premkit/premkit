package persistence

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/premkit/premkit/log"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

// DB is the lazy-loaded reference to the BoltDB instance.  Use the GetDB() function to obtain this.
var DB *bolt.DB
var mu sync.Mutex

// GetDB returns the singleton instance of the BoltDB connection.  This is not a threadsafe object,
// but transactions are.  Any caller using this object should use a transaction.
func GetDB() (*bolt.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if DB != nil {
		return DB, nil
	}

	log.Debugf("Creating connection to data_file at %s", viper.GetString("data_file"))

	if err := os.MkdirAll(filepath.Dir(viper.GetString("data_file")), 0755); err != nil {
		log.Error(err)
		return nil, err
	}

	conn, err := bolt.Open(viper.GetString("data_file"), 0600, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	DB = conn
	return DB, nil
}
