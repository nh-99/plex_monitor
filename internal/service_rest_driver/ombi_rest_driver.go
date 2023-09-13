package servicerestdriver

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// OmbiRestDriver is a REST driver to interact with Ombi. It handles making all standard HTTP requests
type OmbiRestDriver struct {
	ServiceRestDriver
}

// NewOmbiRestDriver returns a new Ombi rest driver.
func NewOmbiRestDriver(name, host, key string, logger *logrus.Entry) *OmbiRestDriver {
	return &OmbiRestDriver{
		ServiceRestDriver: ServiceRestDriver{
			Name:    name,
			Host:    host,
			Key:     key,
			Client:  &http.Client{},
			Logger:  logger,
			Retries: 10,
			Backoff: 5 * time.Second,
		},
	}
}

// GetServiceName returns the name of the service.
func (s *OmbiRestDriver) GetServiceName() string {
	return s.Name
}

// Do executes a request against Ombi.
func (s *OmbiRestDriver) Do(req *http.Request) (*http.Response, error) {
	// Check if the header is nil
	if req.Header == nil {
		// If it is, create a new header
		req.Header = http.Header{}
	}

	// Set the API key header (Ombi requires this convention)
	req.Header.Set("ApiKey", s.Key)

	// Execute the request with the service rest driver
	return s.ExecuteRequestSafe(req)
}
