package persistence

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDB(t *testing.T) {

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
