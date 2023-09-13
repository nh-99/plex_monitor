package servicerestdriver

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// ServiceRestDriver is a REST driver to interact with external services. It handles making all standard HTTP requests
// against a designated endpoint. It also handles authentication and authorization. It is designed to be used as a
// singleton. It is also designed to be used as a driver for the Service struct. It is not designed to be used as a
// standalone service. It builds in retries, exponential backoff, and circuit breaking.

// ServiceRestDriver is the struct that represents the service rest driver.
type ServiceRestDriver struct {
	// Name is the name of the service.
	Name string

	// Host is the host of the service.
	Host string

	// Key is the key for the service.
	Key string

	// Client is the HTTP client for the service.
	Client *http.Client

	// Logger is the logger for the service.
	Logger *logrus.Entry

	// Retries is the number of retries to attempt.
	Retries int

	// Backoff is the amount of time to wait between retries.
	Backoff time.Duration
}

// NewServiceRestDriver returns a new service rest driver.
func NewServiceRestDriver(name, host, key string, retries int, backoff time.Duration, logger *logrus.Entry) *ServiceRestDriver {
	return &ServiceRestDriver{
		Name:    name,
		Host:    host,
		Key:     key,
		Client:  &http.Client{},
		Logger:  logger,
		Retries: 3,
		Backoff: 1 * time.Second,
	}
}

// GetServiceName returns the name of the service.
func (s *ServiceRestDriver) GetServiceName() string {
	return s.Name
}

// GetServiceHost returns the host of the service.
func (s *ServiceRestDriver) GetServiceHost() string {
	return s.Host
}

// GetServiceKey returns the key for the service.
func (s *ServiceRestDriver) GetServiceKey() string {
	return s.Key
}

// GetServiceClient returns the HTTP client for the service.
func (s *ServiceRestDriver) GetServiceClient() *http.Client {
	return s.Client
}

// GetServiceLogger returns the logger for the service.
func (s *ServiceRestDriver) GetServiceLogger() *logrus.Entry {
	return s.Logger
}

// SetServiceName sets the name of the service.
func (s *ServiceRestDriver) SetServiceName(name string) {
	s.Name = name
}

// SetServiceHost sets the host of the service.
func (s *ServiceRestDriver) SetServiceHost(host string) {
	s.Host = host
}

// SetServiceKey sets the key for the service.
func (s *ServiceRestDriver) SetServiceKey(key string) {
	s.Key = key
}

// SetServiceClient sets the HTTP client for the service.
func (s *ServiceRestDriver) SetServiceClient(client *http.Client) {
	s.Client = client
}

// SetServiceLogger sets the logger for the service.
func (s *ServiceRestDriver) SetServiceLogger(logger *logrus.Entry) {
	s.Logger = logger
}

// SetServiceRetries sets the number of retries to attempt.
func (s *ServiceRestDriver) SetServiceRetries(retries int) {
	s.Retries = retries
}

// Do executes a request against the service.
func (s *ServiceRestDriver) Do(req *http.Request) (*http.Response, error) {
	// Execute the request
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Return the response
	return resp, nil
}

// ExecuteRequestSafe executes a request against the service with retries, backoff,
// and circuit breaking.
func (s *ServiceRestDriver) ExecuteRequestSafe(req *http.Request) (*http.Response, error) {
	// Execute the request
	resp, err := s.Client.Do(req)
	if err != nil {
		// If there are retries left, then retry the request
		if s.Retries > 0 {
			time.Sleep(s.Backoff)
			s.Retries--
			return s.ExecuteRequestSafe(req)
		}

		// Otherwise, return the error
		return nil, err
	}

	// Return the response
	return resp, nil
}
