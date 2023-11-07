package worker

import (
	"plex_monitor/internal/database/models"
	servicerestdriver "plex_monitor/internal/service_rest_driver"
	"sort"
	"sync"
	"time"
)

func init() {
	// Register the health cron worker
	RegisterCronWorker(&HealthCronWorker{})
}

// HealthCronWorker is the struct that represents the health cron worker.
type HealthCronWorker struct {
	CronWorker
	latestHealth      map[string]servicerestdriver.ServiceHealth
	latestHealthMutex sync.Mutex
}

// NewHealthCronWorker returns a new health cron worker.
func NewHealthCronWorker() *HealthCronWorker {
	return &HealthCronWorker{}
}

// Name returns the name of the worker.
func (w *HealthCronWorker) Name() string {
	return "Service Health Cron Worker"
}

// GetInterval returns the interval for the worker.
func (w *HealthCronWorker) GetInterval() time.Duration {
	return 10 * time.Minute
}

// Do executes the worker.
func (w *HealthCronWorker) Do() error {
	// Get the logger
	logger := w.GetLogger()
	logger = logger.WithField("worker", w.Name())

	// Log that we are running the health check
	logger.Info("Running health check")

	// Get the services
	services, err := models.GetAllServices()
	if err != nil {
		// Log that we failed to get the services
		logger.Errorf("Failed to get services: %v", err)

		// Return the error
		return err
	}

	// Loop through the services
	for _, service := range services {
		logger = logger.WithField("service", service.ServiceName)
		logger.Debugf("Getting health for service %s", service.ServiceName)

		// Get the driver
		driver, err := servicerestdriver.GetDriverForService(&service, true)
		if err != nil {
			logger.Errorf("Failed to get driver for service: %v", err)
		}

		// Get the health
		health, err := driver.GetHealth()
		if err != nil {
			logger.Errorf("Failed to get health: %v", err)
			// Service is unhealthy - set the health to unhealthy
			health = servicerestdriver.ServiceHealth{
				Healthy: false,
				Version: "N/A",
			}
		}

		// Add the latest health
		w.addLatestHealth(health, string(service.ServiceName))

		logger.Debugf("Service health: %v", health)
	}

	// Log that we are done running the health check
	logger.Info("Done running health check")

	// Return nil
	return nil
}

// GetLatestHealthMap returns the map of latest health checks.
func (w *HealthCronWorker) GetLatestHealthMap() map[string]servicerestdriver.ServiceHealth {
	keys := make([]string, 0, len(w.latestHealth))

	for k := range w.latestHealth {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tempMap := make(map[string]servicerestdriver.ServiceHealth)
	for _, k := range keys {
		tempMap[k] = w.latestHealth[k]
	}

	return tempMap
}

func (w *HealthCronWorker) addLatestHealth(health servicerestdriver.ServiceHealth, serviceName string) {
	// Lock the mutex
	w.latestHealthMutex.Lock()
	defer w.latestHealthMutex.Unlock()

	// Get the logger
	logger := w.GetLogger()

	// Log that we are adding the latest health
	logger.Infof("Adding latest health for service %s", serviceName)

	if w.latestHealth == nil {
		w.latestHealth = make(map[string]servicerestdriver.ServiceHealth)
	}

	// Add the latest health
	w.latestHealth[serviceName] = health
}
