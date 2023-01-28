package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/input"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/metric"
)

func main() {
	// w, err := newWorker()
	// if err != nil {
	// 	panic(err)
	// }
	// w.exec()
}

type worker struct {
	units        int
	shutdown     chan os.Signal
	stopUnits    chan bool
	stoppedUnits chan bool
	inputRepo    input.Repo
	metricRepo   metric.Repo
}

// func newWorker() (*worker, error) {
// 	w := &worker{
// 		units:        cfg.WorkerUnits,
// 		shutdown:     make(chan os.Signal, 1),
// 		stopUnits:    make(chan bool, cfg.WorkerUnits),
// 		stoppedUnits: make(chan bool, cfg.WorkerUnits),
// 	}
// 	return w, w.initModules(cfg.SQS)
// }

// func (w *worker) initModules(sqsURL string) error {
// 	s, err := aws.Session()
// 	if err != nil {
// 		return err
// 	}
// 	d, err := dba.GetDB()
// 	if err != nil {
// 		return err
// 	}
// 	w.inputRepo = repo.NewSQSInputRepo(sqsURL, aws.SQS(s))
// 	w.metricRepo = repo.NewSQLMetricRepo(d)
// 	return nil
// }

func (w *worker) exec() {
	w.initUnits()
	w.handleShutdown()
}

func (w *worker) initUnits() {
	for u := 1; u <= w.units; u++ {
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

	for u := 1; u <= w.units; u++ {
		w.stopUnits <- true
	}

	sppd := 0
	for <-w.stopUnits {
		sppd++
		if sppd == w.units {
			break
		}
	}
	os.Exit(0)
}
