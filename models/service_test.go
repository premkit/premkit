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

func setup(t *testing.T) string {
	dirName, err := ioutil.TempDir("", "premkit-test")
	require.NoError(t, err)

	conn, err := bolt.Open(path.Join(dirName, "test.db"), 0600, nil)
	require.NoError(t, err)

	persistence.DB = conn

	return dirName
}

func teardown(dbPath string) {
	os.RemoveAll(dbPath)
}

func TestListServicesEmpty(t *testing.T) {
	dbPath := setup(t)
	defer teardown(dbPath)

	services, err := ListServices()
	require.NoError(t, err)
	assert.Equal(t, 0, len(services), "there should be no services")
}

func TestCreateService(t *testing.T) {
	dbPath := setup(t)
	defer teardown(dbPath)

	service, err := CreateService(&Service{
		Name: "name",
		Path: "path_a",
		Upstreams: []*Upstream{
			&Upstream{URL: "a"},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "name", service.Name, "service name should be 'name'")
	assert.Equal(t, "path_a", service.Path, "service path should be 'path'")
	assert.Equal(t, 1, len(service.Upstreams), "there should be 1 upstream")
	assert.Equal(t, "a", service.Upstreams[0].URL, "upstream service[0].url should be 'a'")

	services, err := ListServices()
	require.NoError(t, err)
	assert.Equal(t, 1, len(services), "there should be 1 service")
}

func TestUpdateService(t *testing.T) {
	dbPath := setup(t)
	defer teardown(dbPath)

	service, err := CreateService(&Service{
		Name: "name_2",
		Path: "path_2",
		Upstreams: []*Upstream{
			&Upstream{URL: "1"},
		},
	})
	require.NoError(t, err)

	service, err = CreateService(&Service{
		Name: "name_2",
		Path: "path_2",
		Upstreams: []*Upstream{
			&Upstream{URL: "2"},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, "name_2", service.Name, "service name should be 'name'")
	assert.Equal(t, "path_2", service.Path, "service path should be 'path'")
	assert.Equal(t, 2, len(service.Upstreams), "there should be 2 upstreams")
	assert.Equal(t, "1", service.Upstreams[0].URL, "upstream service[0] should be '1'")
	assert.Equal(t, "2", service.Upstreams[1].URL, "upstream service[1] should be '2'")

	services, err := ListServices()
	require.NoError(t, err)
	assert.Equal(t, 1, len(services), "there should be 1 service1")

	var createdService *Service
	for _, s := range services {
		if s.Name == "name_2" {
			createdService = s
		}
	}
	assert.NotNil(t, createdService)
	assert.Equal(t, "name_2", createdService.Name)
	assert.Equal(t, "path_2", createdService.Path)
	assert.Equal(t, 2, len(createdService.Upstreams))
}

func TestDeleteServiceByName(t *testing.T) {
	dbPath := setup(t)
	defer teardown(dbPath)

	deleted, err := DeleteServiceByName([]byte("test"))
	require.NoError(t, err)
	assert.Equal(t, false, deleted)

	_, err = CreateService(&Service{
		Name: "test",
		Path: "path",
		Upstreams: []*Upstream{
			&Upstream{URL: "upstream"},
		},
	})
	require.NoError(t, err)

	deleted, err = DeleteServiceByName([]byte("test"))
	require.NoError(t, err)
	assert.Equal(t, true, deleted)
}
