package servicerestdriver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// PlexRestDriver is a REST driver to interact with Plex. It handles making all standard HTTP requests
type PlexRestDriver struct {
	ServiceRestDriver
}

// NewPlexRestDriver returns a new Plex rest driver.
func NewPlexRestDriver(name, host, key string, logger *logrus.Entry) *PlexRestDriver {
	return &PlexRestDriver{
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
func (s *PlexRestDriver) GetServiceName() string {
	return s.Name
}

// Do executes a request against Plex.
func (s *PlexRestDriver) Do(req *http.Request) (*http.Response, error) {
	// Set the X-Plex-Token in the URL query
	params := req.URL.Query()
	params.Set("X-Plex-Token", s.Key)
	req.URL.RawQuery = params.Encode()

	// Execute the request with the service rest driver
	return s.ExecuteRequestSafe(req)
}

// ScanLibrary scans the library with the given ID.
func (s *PlexRestDriver) ScanLibrary(libraryID int) error {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, s.Host+"/library/sections/"+strconv.Itoa(libraryID)+"/refresh", nil)
	if err != nil {
		return err
	}

	// Execute the request
	_, err = s.Do(req)
	if err != nil {
		return err
	}

	return nil
}
