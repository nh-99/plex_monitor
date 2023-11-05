package servicerestdriver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// NotificationPreferenceDiscordAgentID is the ID of the Discord notification preference
	NotificationPreferenceDiscordAgentID = 1
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

// OmbiUser is the struct that represents a user in Ombi
type OmbiUser struct {
	ID           string `json:"id"`
	UserName     string `json:"userName"`
	Alias        string `json:"alias"`
	EmailAddress string `json:"emailAddress"`
}

// NotificationPreference is the struct that represents a notification preference in Ombi
type NotificationPreference struct {
	UserID  string  `json:"userId"`
	AgentID int     `json:"agent"`
	Enabled bool    `json:"enabled"`
	Value   *string `json:"value,omitempty"`
	ID      int     `json:"id"`
}

// GetServerVersion returns the version of the Ombi server.
func (s *OmbiRestDriver) GetServerVersion() (string, error) {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, s.Host+"/api/v1/Status/info", nil)
	if err != nil {
		return "", err
	}

	// Execute the request
	resp, err := s.Do(req)
	if err != nil {
		return "", err
	}

	// The response comes back as just a string of the version, so convert the whole body to a string
	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	versionStr := string(reader)

	// Remove quotes from the version string
	versionStr = versionStr[1 : len(versionStr)-1]

	// Return the version
	return versionStr, nil
}

// RetrieveUsersFromOmbi returns the list of users from Ombi
func (s *OmbiRestDriver) RetrieveUsersFromOmbi() ([]OmbiUser, error) {
	// Reach out to the media request service and get the user's Discord ID.
	usersListURL := fmt.Sprintf("%s/ombi/api/v1/Identity/Users", s.Host)
	getUsersListRequest, err := http.NewRequest(http.MethodGet, usersListURL, nil)
	resp, err := s.Do(getUsersListRequest)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response into a list of users.
	var users []OmbiUser
	err = json.Unmarshal(body, &users)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return users, nil
}

// RetrieveDiscordIDFromOmbi returns the Discord ID of the user who requested the media
func (s *OmbiRestDriver) RetrieveDiscordIDFromOmbi(user OmbiUser) (string, error) {
	// Reach out to the media request service and get the user's Discord ID.
	userNotificationPrefsURL := fmt.Sprintf("%s/ombi/api/v1/Identity/notificationpreferences/%s", s.Host, user.ID)
	getUserNotificationPrefsRequest, err := http.NewRequest(http.MethodGet, userNotificationPrefsURL, nil)
	if err != nil {
		s.Logger.WithFields(logrus.Fields{
			"url":          userNotificationPrefsURL,
			"ombiUserName": user.UserName,
			"ombiUserId":   user.ID,
			"error":        err,
		}).Error(err)
		return "", err
	}
	resp, err := s.Do(getUserNotificationPrefsRequest)
	if err != nil {
		s.Logger.WithFields(logrus.Fields{
			"url":          userNotificationPrefsURL,
			"ombiUserName": user.UserName,
			"ombiUserId":   user.ID,
			"error":        err,
		}).Error(err)
		return "", err
	}

	// Read the response body.
	responseBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		s.Logger.WithFields(logrus.Fields{
			"responseBody": string(responseBody),
			"ombiUserName": user.UserName,
			"ombiUserId":   user.ID,
			"error":        err,
		}).Error(err)
		return "", err
	}

	// Parse the response into a list of notification preferences.
	var notificationPreferences []NotificationPreference
	err = json.Unmarshal(responseBody, &notificationPreferences)
	if err != nil {
		s.Logger.WithFields(logrus.Fields{
			"notificationPreferencesCount": len(notificationPreferences),
			"ombiUserName":                 user.UserName,
			"ombiUserId":                   user.ID,
			"error":                        err,
		}).Error(err)
		return "", err
	}

	// Filter the notification preferences to the agent ID for Discord.
	for _, notificationPreference := range notificationPreferences {
		if notificationPreference.AgentID == NotificationPreferenceDiscordAgentID && notificationPreference.UserID == user.ID {
			return *notificationPreference.Value, nil
		}
	}

	// If we get here, then we couldn't find the Discord ID for the user.
	return "", fmt.Errorf("could not find Discord ID for user %s", user.UserName)
}

// GetHealth returns the health of the Ombi server.
func (s *OmbiRestDriver) GetHealth() (ServiceHealth, error) {
	// Get the server capabilities
	version, err := s.GetServerVersion()
	if err != nil {
		return ServiceHealth{
			Healthy:     false,
			Version:     "N/A",
			LastChecked: time.Now(),
		}, fmt.Errorf("failed to get server version: %w", err)
	}

	// Create the health struct
	health := ServiceHealth{
		Healthy:     true,
		Version:     version,
		LastChecked: time.Now(),
	}

	return health, nil
}
