package persistence

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDB(t *testing.T) {
	// Ensure we start with no data dir
	if DB != nil {
		DB.Close()
		DB = nil
	}

	originalDataFile := viper.GetString("data_file")
	defer viper.Set("data_file", originalDataFile)

	viper.Set("data_file", "/tmp/testing.db")

	db, err := GetDB()
	require.NoError(t, err)
	assert.NotNil(t, db)

	db2, err := GetDB()
	require.NoError(t, err)
	assert.Equal(t, db, db2)

	db.Close()
	os.RemoveAll("/tmp/testing.db")
}

func TestInitializeDatabase(t *testing.T) {
	dirName, err := ioutil.TempDir("", "premkit-test")
	require.NoError(t, err)

	conn, err := bolt.Open(path.Join(dirName, "test.db"), 0600, nil)
	require.NoError(t, err)

	err = initializeDatabase(conn)
	require.NoError(t, err)

	itemCount := 0
	err = conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Services"))

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			itemCount++
		}

		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, 0, itemCount, "there should be 0 items in the Services bucket immediately after initialization")
}
