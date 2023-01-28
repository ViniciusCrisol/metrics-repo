package mysql

import (
	"database/sql"
	"time"

	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/internal/config"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/log"
)

const (
	kind            = "mysql"
	maxIdleConns    = 0
	maxOpenConns    = 100
	connMaxLifetime = time.Minute
)

func NewConn() (*sql.DB, error) {
	c, err := sql.Open(kind, config.DBConnURL)
	if err != nil {
		log.Logger.Error(
			"Failed to init DB session",
			log.Error(err),
		)
		return nil, err
	}
	c.SetMaxIdleConns(maxIdleConns)
	c.SetMaxOpenConns(maxOpenConns)
	c.SetConnMaxLifetime(connMaxLifetime)

	if err = c.Ping(); err != nil {
		log.Logger.Error(
			"Failed to establish DB connection",
			log.Error(err),
		)
		return nil, err
	}
	return c, nil
}
