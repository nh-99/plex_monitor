package worker

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// CronWorkers is a list of all cron workers.
	CronWorkers    []CronWorkerImpl
	cronWorkerOnce sync.Once
)

// CronWorker is the struct that represents a cron worker.
type CronWorker struct {
	logger *logrus.Entry
}

// CronWorkerImpl is the interface that represents a cron worker.
type CronWorkerImpl interface {
	// Name returns the name of the worker.
	Name() string

	// Do executes the worker.
	Do() error

	// GetInterval returns the interval for the worker.
	GetInterval() time.Duration

	// GetLogger returns the logger for the worker.
	GetLogger() *logrus.Entry

	// SetLogger sets the logger for the worker.
	SetLogger(logger *logrus.Entry)
}

// RegisterCronWorker registers a cron worker.
func RegisterCronWorker(worker CronWorkerImpl) {
	worker.SetLogger(logrus.WithField("worker", worker.Name()))
	// Append the cron worker to the list of cron workers
	CronWorkers = append(CronWorkers, worker)
}

// Name returns the name of the worker.
func (w *CronWorker) Name() string {
	return "Cron Worker"
}

// GetLogger returns the logger for the worker.
func (w *CronWorker) GetLogger() *logrus.Entry {
	return w.logger
}

// SetLogger sets the logger for the worker.
func (w *CronWorker) SetLogger(logger *logrus.Entry) {
	w.logger = logger
}

// ExecuteCrons executes all cron workers.
func ExecuteCrons() {
	// Run all sync first to ensure that all data is up to date
	ExecuteAllCronWorkersSync()
	// Execute the cron workers
	ExecuteCronWorkers()
}

// ExecuteCronWorkers executes all cron workers.
func ExecuteCronWorkers() {
	logrus.Info("Executing cron workers")
	cronWorkerOnce.Do(func() {
		// Loop through the cron workers
		for _, worker := range CronWorkers {
			logrus.Infof("Executing cron worker: %s", worker.Name())
			ticker := time.NewTicker(worker.GetInterval())
			quit := make(chan struct{})
			go func(worker CronWorkerImpl) {
				for {
					select {
					case <-ticker.C:
						// do stuff
						err := worker.Do()
						if err != nil {
							logrus.Errorf("Failed to execute cron worker: %v", err)
						}
					case <-quit:
						ticker.Stop()
						return
					}
				}
			}(worker)
		}
	})
}

// ExecuteAllCronWorkersSync executes all cron workers synchronously.
func ExecuteAllCronWorkersSync() {
	// Loop through the cron workers
	for _, worker := range CronWorkers {
		// Execute the worker
		err := worker.Do()
		if err != nil {
			logrus.Errorf("Failed to execute cron worker: %v", err)
		}
	}
}

// GetCronWorkers returns all cron workers.
func GetCronWorkers() []CronWorkerImpl {
	return CronWorkers
}

// GetCronWorker returns the cron worker with the given name.
func GetCronWorker(name string) CronWorkerImpl {
	// Loop through the cron workers
	for _, worker := range CronWorkers {
		// Check if the names match
		if worker.Name() == name {
			return worker
		}
	}

	return nil
}
