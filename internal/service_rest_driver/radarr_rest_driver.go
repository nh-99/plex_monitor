package servicerestdriver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// RadarrRestDriver is a REST driver to interact with Radarr. It handles making all standard HTTP requests
type RadarrRestDriver struct {
	ServiceRestDriver
}

// NewRadarrRestDriver returns a new Plex rest driver.
func NewRadarrRestDriver(name, host, key string, logger *logrus.Entry) *RadarrRestDriver {
	return &RadarrRestDriver{
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

// Do executes a request against Radarr.
func (s *RadarrRestDriver) Do(req *http.Request) (*http.Response, error) {
	// Check if the header is nil
	if req.Header == nil {
		// If it is, create a new header
		req.Header = http.Header{}
	}

	// Set the API key header (Radarr requires this convention)
	req.Header.Set("X-Api-Key", s.Key)

	// Execute the request with the service rest driver
	return s.ExecuteRequestSafe(req)
}

// RadarrSystemStatus is the struct that represents the system status of Radarr
type RadarrSystemStatus struct {
	AppName                       string    `json:"appName"`
	InstanceName                  string    `json:"instanceName"`
	Version                       string    `json:"version"`
	BuildTime                     time.Time `json:"buildTime"`
	IsDebug                       bool      `json:"isDebug"`
	IsProduction                  bool      `json:"isProduction"`
	IsAdmin                       bool      `json:"isAdmin"`
	IsUserInteractive             bool      `json:"isUserInteractive"`
	StartupPath                   string    `json:"startupPath"`
	AppData                       string    `json:"appData"`
	OsName                        string    `json:"osName"`
	OsVersion                     string    `json:"osVersion"`
	IsNetCore                     bool      `json:"isNetCore"`
	IsLinux                       bool      `json:"isLinux"`
	IsOsx                         bool      `json:"isOsx"`
	IsWindows                     bool      `json:"isWindows"`
	IsDocker                      bool      `json:"isDocker"`
	Mode                          string    `json:"mode"`
	Branch                        string    `json:"branch"`
	DatabaseType                  string    `json:"databaseType"`
	DatabaseVersion               string    `json:"databaseVersion"`
	Authentication                string    `json:"authentication"`
	MigrationVersion              int       `json:"migrationVersion"`
	URLBase                       string    `json:"urlBase"`
	RuntimeVersion                string    `json:"runtimeVersion"`
	RuntimeName                   string    `json:"runtimeName"`
	StartTime                     time.Time `json:"startTime"`
	PackageVersion                string    `json:"packageVersion"`
	PackageAuthor                 string    `json:"packageAuthor"`
	PackageUpdateMechanism        string    `json:"packageUpdateMechanism"`
	PackageUpdateMechanismMessage string    `json:"packageUpdateMechanismMessage"`
}

// GetSystemStatus returns the system status of the service.
func (s *RadarrRestDriver) GetSystemStatus() (RadarrSystemStatus, error) {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, s.Host+"/api/v3/system/status", nil)
	if err != nil {
		return RadarrSystemStatus{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request
	resp, err := s.Do(req)
	if err != nil {
		return RadarrSystemStatus{}, fmt.Errorf("failed to execute request: %w", err)
	}

	// Close the response body
	defer resp.Body.Close()

	// Decode the response
	var status RadarrSystemStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return RadarrSystemStatus{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return status, nil
}

// GetHealth returns the health of the service.
func (s *RadarrRestDriver) GetHealth() (ServiceHealth, error) {
	// Use the system status to get the health
	status, err := s.GetSystemStatus()
	if err != nil {
		return ServiceHealth{
			Healthy:     false,
			Version:     "N/A",
			LastChecked: time.Now(),
		}, err
	}

	// Convert the system status to a service health
	return ServiceHealth{
		Healthy:     true,
		Version:     status.Version,
		LastChecked: time.Now(),
	}, nil
}
