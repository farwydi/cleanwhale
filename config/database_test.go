package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConfig_Load(t *testing.T) {
	var (
		ds     DatabaseConfig
		err    error
		source string
	)

	require.NoError(t, LoadConfigs(&ds, "../testdata/datasource.yml"))

	source, err = ds.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "host=dbhost user=postgres password= dbname=dbname sslmode=disable", source)

	assert.Equal(t, 2, ds.MaxIdleConns)
	assert.Equal(t, 0, ds.MaxOpenConns)
	assert.Equal(t, time.Minute, ds.ConnMaxLifetime)
}

func TestDatabaseConfig_ToString_Errs(t *testing.T) {
	var (
		ds     DatabaseConfig
		err    error
		source string
	)

	ds = DatabaseConfig{}
	_, err = ds.ToString()
	assert.Equal(t, ErrSetupHost, err)

	ds = DatabaseConfig{Host: "dbhost"}
	_, err = ds.ToString()
	assert.Equal(t, ErrSetupName, err)

	ds = DatabaseConfig{Host: "dbhost", Name: "dbname"}
	source, err = ds.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "host=dbhost user= password= dbname=dbname sslmode=disable", source)
}
