package servicerestdriver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// SonarrRestDriver is a REST driver to interact with Sonarr. It handles making all standard HTTP requests
type SonarrRestDriver struct {
	ServiceRestDriver
}

// NewSonarrRestDriver returns a new Plex rest driver.
func NewSonarrRestDriver(name, host, key string, logger *logrus.Entry) *SonarrRestDriver {
	return &SonarrRestDriver{
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

// Do executes a request against Sonarr.
func (s *SonarrRestDriver) Do(req *http.Request) (*http.Response, error) {
	// Check if the header is nil
	if req.Header == nil {
		// If it is, create a new header
		req.Header = http.Header{}
	}

	// Set the API key header (Sonarr requires this convention)
	req.Header.Set("X-Api-Key", s.Key)

	// Execute the request with the service rest driver
	return s.ExecuteRequestSafe(req)
}

// SonarrSystemStatus is the struct that represents the system status of Sonarr
type SonarrSystemStatus struct {
	AppName                       string    `json:"appName" bson:"appName"`
	InstanceName                  string    `json:"instanceName" bson:"instanceName"`
	Version                       string    `json:"version" bson:"version"`
	BuildTime                     time.Time `json:"buildTime" bson:"buildTime"`
	IsDebug                       bool      `json:"isDebug" bson:"isDebug"`
	IsProduction                  bool      `json:"isProduction" bson:"isProduction"`
	IsAdmin                       bool      `json:"isAdmin" bson:"isAdmin"`
	IsUserInteractive             bool      `json:"isUserInteractive" bson:"isUserInteractive"`
	StartupPath                   string    `json:"startupPath" bson:"startupPath"`
	AppData                       string    `json:"appData" bson:"appData"`
	OsName                        string    `json:"osName" bson:"osName"`
	OsVersion                     string    `json:"osVersion" bson:"osVersion"`
	IsNetCore                     bool      `json:"isNetCore" bson:"isNetCore"`
	IsLinux                       bool      `json:"isLinux" bson:"isLinux"`
	IsOsx                         bool      `json:"isOsx" bson:"isOsx"`
	IsWindows                     bool      `json:"isWindows" bson:"isWindows"`
	IsDocker                      bool      `json:"isDocker" bson:"isDocker"`
	Mode                          string    `json:"mode" bson:"mode"`
	Branch                        string    `json:"branch" bson:"branch"`
	Authentication                string    `json:"authentication" bson:"authentication"`
	SqliteVersion                 string    `json:"sqliteVersion" bson:"sqliteVersion"`
	MigrationVersion              int       `json:"migrationVersion" bson:"migrationVersion"`
	URLBase                       string    `json:"urlBase" bson:"urlBase"`
	RuntimeVersion                string    `json:"runtimeVersion" bson:"runtimeVersion"`
	RuntimeName                   string    `json:"runtimeName" bson:"runtimeName"`
	StartTime                     time.Time `json:"startTime" bson:"startTime"`
	PackageVersion                string    `json:"packageVersion" bson:"packageVersion"`
	PackageAuthor                 string    `json:"packageAuthor" bson:"packageAuthor"`
	PackageUpdateMechanism        string    `json:"packageUpdateMechanism" bson:"packageUpdateMechanism"`
	PackageUpdateMechanismMessage string    `json:"packageUpdateMechanismMessage" bson:"packageUpdateMechanismMessage"`
	DatabaseVersion               string    `json:"databaseVersion" bson:"databaseVersion"`
	DatabaseType                  string    `json:"databaseType" bson:"databaseType"`
}

// GetSystemStatus returns the system status of the service.
func (s *SonarrRestDriver) GetSystemStatus() (SonarrSystemStatus, error) {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, s.Host+"/api/v3/system/status", nil)
	if err != nil {
		return SonarrSystemStatus{}, err
	}

	// Execute the request
	resp, err := s.Do(req)
	if err != nil {
		return SonarrSystemStatus{}, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Decode the response
	var status SonarrSystemStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return SonarrSystemStatus{}, err
	}

	return status, nil
}

// GetHealth returns the health of the service.
func (s *SonarrRestDriver) GetHealth() (ServiceHealth, error) {
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
