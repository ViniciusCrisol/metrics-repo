package mysql

import (
	"database/sql"
	"time"

	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/internal/config"
)

const (
	kind            = "mysql"
	maxIdleConns    = 0
	maxOpenConns    = 100
	connMaxLifetime = time.Minute
)

func NewDBConn() (*sql.DB, error) {
	c, err := sql.Open(kind, config.DBConnURL)
	if err != nil {
		// TODO: Log it!
		return nil, err
	}
	c.SetMaxIdleConns(maxIdleConns)
	c.SetMaxOpenConns(maxOpenConns)
	c.SetConnMaxLifetime(connMaxLifetime)

	if err = c.Ping(); err != nil {
		// TODO: Log it!
		return nil, err
	}
	return c, nil
}
