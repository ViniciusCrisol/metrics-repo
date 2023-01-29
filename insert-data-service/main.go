package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/internal/config"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/internal/config/aws"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/internal/config/mysql"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/internal/repo"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/input"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/metric"
)

const workUnits = 5

func main() {
	w, err := newWorker()
	if err != nil {
		panic(err)
	}
	w.Exec()
}

type worker struct {
	shutdown     chan os.Signal
	stopUnits    chan bool
	stoppedUnits chan bool

	inputRepo  input.Repo
	metricRepo metric.Repo
}

func newWorker() (*worker, error) {
	w := &worker{
		shutdown:     make(chan os.Signal, 1),
		stopUnits:    make(chan bool, workUnits),
		stoppedUnits: make(chan bool, workUnits),
	}
	return w, w.initModules(config.SQS)
}

func (w *worker) initModules(sqsURL string) error {
	s, err := aws.NewSession()
	if err != nil {
		return err
	}
	d, err := mysql.NewConn()
	if err != nil {
		return err
	}
	w.metricRepo = repo.NewMetricSQLRepo(d)
	w.inputRepo = repo.NewInputSQSRepo(
		sqsURL,
		aws.NewSQS(s),
	)
	return nil
}

func (w *worker) Exec() {
	w.initUnits()
	w.handleShutdown()
}

func (w *worker) initUnits() {
	for u := 1; u <= workUnits; u++ {
		go w.initUnit()
	}
}

func (w *worker) initUnit() {
	for {
		select {
		case <-w.stopUnits:
			w.stopUnit()
			return
		default:
			is, err := w.inputRepo.Get()
			if err != nil {
				continue
			}
			if len(is) == 0 {
				time.Sleep(time.Minute)
				continue
			}
			wg := &sync.WaitGroup{}

			for _, i := range is {
				wg.Add(1)
				w.processInput(wg, i)
			}
		}
	}
}

func (w *worker) stopUnit() {
	w.stoppedUnits <- true
}

func (w *worker) processInput(wg *sync.WaitGroup, i *input.Input) {
	defer wg.Done()

	m, err := metric.FromInput(i)
	if err != nil {
		return
	}
	if err = w.metricRepo.Create(m); err != nil {
		return
	}
	w.inputRepo.Delete(i)
}

func (w *worker) handleShutdown() {
	signal.Notify(w.shutdown, syscall.SIGTERM)

	<-w.shutdown

	for u := 1; u <= workUnits; u++ {
		w.stopUnits <- true
	}

	sppd := 0
	for <-w.stopUnits {
		sppd++
		if sppd == workUnits {
			break
		}
	}
	os.Exit(0)
}
