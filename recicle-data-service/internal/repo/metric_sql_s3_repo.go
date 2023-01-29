package repo

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/service/s3"
)

type metricSQLS3Repo struct {
	s3     *s3.S3
	dbConn *sql.DB
}

func NewMetricSQLS3Repo(s3 *s3.S3, dbConn *sql.DB) *metricSQLS3Repo {
	return &metricSQLS3Repo{
		s3:     s3,
		dbConn: dbConn,
	}
}

func (repo *metricSQLS3Repo) Recicle() error {
	// TODO: Implement it!
	return nil
}
