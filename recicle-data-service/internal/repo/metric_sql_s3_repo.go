package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ViniciusCrisol/metrics-repo/recicle-data-service/log"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type metricSQLS3Repo struct {
	dbConn   *sql.DB
	bucket   string
	uploader *s3manager.Uploader
}

func NewMetricSQLS3Repo(dbConn *sql.DB, bucket string, uploader *s3manager.Uploader) *metricSQLS3Repo {
	return &metricSQLS3Repo{
		dbConn:   dbConn,
		bucket:   bucket,
		uploader: uploader,
	}
}

func (repo *metricSQLS3Repo) Recicle() error {
	m, err := repo.getExpiredMetrics()
	if err != nil {
		return err
	}
	if len(m) == 0 {
		return nil
	}
	return repo.sendToS3AndDeleteFromDB(m)
}

type metric struct {
	ID      int
	Data    string
	AppName string
}

func (repo *metricSQLS3Repo) getExpiredMetrics() ([]*metric, error) {
	rows, err := repo.dbConn.Query(`
		select
			id,
			data,
			app_name
		from
			app_metrics
		where
			created_at < NOW() - INTERVAL 1 WEEK
	`)
	if err != nil {
		log.Logger.Error("Failed to get expired metric", log.Error(err))
		return nil, err
	}
	return repo.rowsToMetrics(rows)
}

func (repo *metricSQLS3Repo) rowsToMetrics(rows *sql.Rows) ([]*metric, error) {
	ms := []*metric{}

	for rows.Next() {
		m := metric{}

		err := rows.Scan(
			&m.ID,
			&m.Data,
			&m.AppName,
		)
		if err != nil {
			log.Logger.Error("Failed to parse row to metric", log.Error(err))
			return nil, err
		}
		ms = append(ms, &m)
	}
	return ms, nil
}

func (repo *metricSQLS3Repo) sendToS3AndDeleteFromDB(ms []*metric) error {
	tx, err := repo.dbConn.Begin()
	if err != nil {
		log.Logger.Error("Failed to begin", log.Error(err))
		return err
	}
	defer func() {
		defer tx.Rollback()

		if p := recover(); p != nil {
			panic(p)
		}
	}()

	if err = repo.deleteFromDB(tx, ms); err != nil {
		return err
	}
	if err = repo.sendToS3(ms); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		log.Logger.Error("Failed to commit", log.Error(err))
		return err
	}
	return nil
}

func (repo *metricSQLS3Repo) deleteFromDB(tx *sql.Tx, ms []*metric) error {
	_, err := tx.Exec(
		repo.getDeleteCommand(ms),
	)
	if err != nil {
		log.Logger.Error("Failed to delete metrics", log.Error(err))
		return err
	}
	return nil
}

func (repo *metricSQLS3Repo) getDeleteCommand(ms []*metric) string {
	ids := ""

	for _, m := range ms {
		ids += fmt.Sprintf("%d,", m.ID)
	}
	ids = strings.TrimSuffix(ids, ",")

	return fmt.Sprintf("delete from app_metrics where id in (%s)", ids)
}

func (repo *metricSQLS3Repo) sendToS3(ms []*metric) error {
	m, err := json.Marshal(ms)
	if err != nil {
		log.Logger.Error("Failed to marshal metrics", log.Error(err))
		return err
	}
	fileName := fmt.Sprintf("%d", time.Now().Unix())

	_, err = repo.uploader.Upload(
		&s3manager.UploadInput{
			Body: strings.NewReader(
				string(m),
			),
			Key:    &fileName,
			Bucket: &repo.bucket,
		},
	)
	if err != nil {
		log.Logger.Error("Failed to upload files", log.Error(err))
		return err
	}
	return nil
}
