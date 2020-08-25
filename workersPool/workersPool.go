package workerspool

import (
	"runtime"
	"time"

	"github.com/aanufriev/SoftProTest/storage"
	"github.com/sirupsen/logrus"
)

type WorkersPool struct {
	workersCount int
	Sports       []string
	Intervals    []int
	workChan     chan string
	storage      storage.DatabaseInterface
}

func NewWorkersPool(workersCount int, sports []string, intervals []int, storage storage.DatabaseInterface) WorkersPool {
	return WorkersPool{
		workersCount: workersCount,
		Sports:       sports,
		Intervals:    intervals,
		workChan:     make(chan string, 10),
		storage:      storage,
	}
}

func (pool WorkersPool) Start(linesProviderURL string) {
	pool.AddWork()
	for i := 0; i < pool.workersCount; i++ {
		worker := Worker{
			url:     linesProviderURL,
			storage: pool.storage,
		}

		go func(worker Worker) {
			for {
				sport := <-pool.workChan
				worker.GetLine(sport)
				runtime.Gosched()
			}
		}(worker)
	}

	logrus.WithFields(logrus.Fields{
		"workers count":      pool.workersCount,
		"lines provider URL": linesProviderURL,
	}).Info("Workers started")
}

func (pool WorkersPool) AddWork() {
	for i, sport := range pool.Sports {
		interval := pool.Intervals[i]
		go func(sport string, interval int) {
			for {
				pool.workChan <- sport
				time.Sleep(time.Second * time.Duration(interval))
			}
		}(sport, interval)
	}
}
