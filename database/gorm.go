// Package database defines auxiliary functions for working with databases.
package database

import (
	"database/sql"
	"time"

	"github.com/farwydi/cleanwhale/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDatabase create gorm.DB instance all log goes to zlog logger.
// Also returns sql.DB for compatibility, if you do not intend to use it, you can safely ignore it.
// The logger will be configured depending on the Mode.
func NewGormDatabase(cfg config.DatabaseConfig, mode config.Mode, zlog *zap.Logger) (*gorm.DB, *sql.DB, func(), error) {
	dbConfig := &gorm.Config{
		Logger: logger.New(
			zap.NewStdLog(zlog),
			logger.Config{
				// Slow SQL threshold
				SlowThreshold: time.Second,
				// Log level
				LogLevel: logger.Info,
				// Disable color
				Colorful: true,
			},
		),
	}

	if mode == config.ModeRelease {
		dbConfig.Logger.LogMode(logger.Silent)
	}

	ds, err := cfg.ToString()
	if err != nil {
		return nil, nil, nil, errors.Wrap(
			errors.WithMessage(err, "config to source"), "new gorm database")
	}

	db, err := gorm.Open(postgres.Open(ds), dbConfig)
	if err != nil {
		return nil, nil, nil, errors.Wrap(
			errors.WithMessage(err, "gorm Open (postgres)"), "new gorm database")
	}

	dd, err := db.DB()
	if err != nil {
		return nil, nil, nil, errors.Wrap(
			errors.WithMessage(err, "take sql.DB"), "new gorm database")
	}

	cleanup := func() {
		if err := dd.Close(); err != nil {
			zlog.Error("Fail close database content",
				zap.String("link", ds),
				zap.Error(err),
			)
		}
	}

	dd.SetMaxOpenConns(cfg.MaxOpenConns)
	dd.SetMaxIdleConns(cfg.MaxIdleConns)
	dd.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, dd, cleanup, nil
}
