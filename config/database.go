package config

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrSetupHost emerge if database host is not setup
	ErrSetupHost = errors.New("setup host in database config")
	// ErrSetupName emerge if database name is not setup
	ErrSetupName = errors.New("setup name in database config")
)

// The DatabaseConfig struct for define database config.
// Method DatabaseConfig.ToString transform config to DS string.
type DatabaseConfig struct {
	Name     string
	User     string `default:"postgres"`
	Password string `default:""`
	Host     string

	MaxOpenConns    int           `yaml:"max_open_conns" default:"0"`
	MaxIdleConns    int           `yaml:"max_idle_conns" default:"2"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" default:"1m"`
}

// ToString transform DatabaseConfig to DS string.
func (ds DatabaseConfig) ToString() (string, error) {
	if ds.Host == "" {
		return "", ErrSetupHost
	}

	if ds.Name == "" {
		return "", ErrSetupName
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		ds.Host,
		ds.User,
		ds.Password,
		ds.Name,
	), nil
}
