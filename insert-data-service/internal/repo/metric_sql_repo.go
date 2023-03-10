package repo

import (
	"database/sql"
	"time"

	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/log"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/metric"
)

type metricSQLRepo struct {
	dbConn *sql.DB
}

func NewMetricSQLRepo(dbConn *sql.DB) *metricSQLRepo {
	return &metricSQLRepo{dbConn}
}

const insertMetricSQLCommand = `
	insert into app_metrics(data, app_name, created_at)
	values(?, ?, ?);
`

func (repo *metricSQLRepo) Create(m *metric.Metric) error {
	_, err := repo.dbConn.Exec(
		insertMetricSQLCommand,
		m.Data,
		m.AppName,
		time.Now(),
	)
	if err != nil {
		log.Logger.Error(
			"Failed to insert a new metric",
			log.Error(err),
			log.String("data", m.Data),
			log.String("app_name", m.AppName),
		)
		return err
	}
	return nil
}
