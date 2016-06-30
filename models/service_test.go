package models

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/premkit/premkit/persistence"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Initialize a temp boltdb folder
	dbPath, err := setup()
	if err != nil {
		os.Exit(1)
	}
	defer teardown(dbPath)

	// Run the tests
	result := m.Run()

	os.Exit(result)
}

func setup() (string, error) {
	dirName, err := ioutil.TempDir("", "premkit-test")
	if err != nil {
		return "", err
	}

	conn, err := bolt.Open(path.Join(dirName, "test.db"), 0600, nil)
	if err != nil {
		return "", err
	}

	err = conn.Update(func(tx *bolt.Tx) error {
		// Create the default buckets
		_, err := tx.CreateBucketIfNotExists([]byte("Services"))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	persistence.DB = conn

	return dirName, nil
}

func teardown(dbPath string) {
	os.RemoveAll(dbPath)
}

func Test_0_ListServicesEmpty(t *testing.T) {
	services, err := ListServices()
	require.NoError(t, err)
	assert.Equal(t, 0, len(services), "there should be no services")
}

func Test_1_CreateService(t *testing.T) {
	service, err := CreateService(&Service{
		Name:      "name",
		Path:      "path",
		Upstreams: []string{"1"},
	})

	require.NoError(t, err)
	assert.Equal(t, "name", service.Name, "service name should be 'name'")
	assert.Equal(t, "path", service.Path, "service path should be 'path'")
	assert.Equal(t, 1, len(service.Upstreams), "there should be 1 upsream")
	assert.Equal(t, "1", service.Upstreams[0], "upstream service[0] should be '1'")

	services, err := ListServices()
	require.NoError(t, err)
	assert.Equal(t, 1, len(services), "there should be 1 service")
}

func Test_2_UpdateService(t *testing.T) {
	service, err := CreateService(&Service{
		Name:      "name",
		Path:      "path",
		Upstreams: []string{"1"},
	})
	require.NoError(t, err)

	service, err = CreateService(&Service{
		Name:      "name",
		Path:      "path",
		Upstreams: []string{"2"},
	})
	require.NoError(t, err)

	assert.Equal(t, "name", service.Name, "service name should be 'name'")
	assert.Equal(t, "path", service.Path, "service path should be 'path'")
	assert.Equal(t, 2, len(service.Upstreams), "there should be 2 upstreams")
	assert.Equal(t, "1", service.Upstreams[0], "upstream service[0] should be '1'")
	assert.Equal(t, "2", service.Upstreams[1], "upstream service[1] should be '2'")

	services, err := ListServices()
	require.NoError(t, err)
	assert.Equal(t, 1, len(services), "there should be 1 service")
}
