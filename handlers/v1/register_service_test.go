package v1

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/premkit/premkit/models"
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

func TestRegisterService(t *testing.T) {
	dbPath := setup(t)
	defer teardown(dbPath)

	params := RegisterServiceParams{
		ReplaceExisting: false,
		Service: &models.Service{
			Name: "test",
			Path: "test",
			Upstreams: []*models.Upstream{
				&models.Upstream{
					URL: "url",
				},
			},
		},
	}

	service, err := registerService(&params)
	require.NoError(t, err)
	assert.NotNil(t, service)
}
