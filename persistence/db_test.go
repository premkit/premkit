package persistence

import (
	"os"
	"testing"

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
