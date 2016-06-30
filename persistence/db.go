package persistence

import (
	"github.com/premkit/premkit/log"

	"github.com/boltdb/bolt"
)

// DB is the lazy-loaded reference to the BoltDB instance.  Use the GetDB() function to obtain this.
var DB *bolt.DB

// GetDB returns the singleton instance of the BoltDB connection.  This is not a threadsafe object,
// but transactions are.  Any caller using this object should use a transaction.
func GetDB() (*bolt.DB, error) {
	if DB != nil {
		return DB, nil
	}

	conn, err := bolt.Open("/data/premkit.db", 0600, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := initializeDatabase(conn); err != nil {
		return nil, err
	}

	DB = conn
	return DB, nil
}

func initializeDatabase(conn *bolt.DB) error {
	// Perform some initialization
	err := conn.Update(func(tx *bolt.Tx) error {
		// Create the default buckets
		_, err := tx.CreateBucketIfNotExists([]byte("Services"))
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	return err
}
