package main

import (
	"time"

	"github.com/ViniciusCrisol/metrics-repo/recicle-data-service/internal/config"
	"github.com/ViniciusCrisol/metrics-repo/recicle-data-service/internal/config/aws"
	"github.com/ViniciusCrisol/metrics-repo/recicle-data-service/internal/config/mysql"
	"github.com/ViniciusCrisol/metrics-repo/recicle-data-service/internal/repo"
)

const maxNumOfRetries = 5

func main() {
	w, err := newWorker()
	if err != nil {
		panic(err)
	}
	w.Exec()
}

type worker struct {
	retries    int
	metricRepo interface{ Recicle() error }
}

func newWorker() (*worker, error) {
	w := &worker{
		retries: 0,
	}
	return w, w.initModules()
}

func (w *worker) initModules() error {
	s, err := aws.NewSession()
	if err != nil {
		return err
	}
	d, err := mysql.NewConn()
	if err != nil {
		return err
	}
	w.metricRepo = repo.NewMetricSQLS3Repo(
		d,
		config.Bucket,
		aws.NewS3(s),
	)
	return nil
}

func (w *worker) Exec() {
	if err := w.metricRepo.Recicle(); err != nil && w.retries < maxNumOfRetries {
		time.Sleep(time.Minute)
		w.retries++
		w.Exec()
	}
}
