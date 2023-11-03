package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"plex_monitor/internal/database"
	servicerestdriver "plex_monitor/internal/service_rest_driver"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// NotificationPreferenceDiscordAgentID is the ID of the Discord notification preference
	NotificationPreferenceDiscordAgentID = 1
)

// OmbiWebhookData is the struct that represents the data sent by Ombi
type OmbiWebhookData struct {
	RequestID                        string    `json:"requestId" bson:"requestId"`
	RequestedUser                    string    `json:"requestedUser" bson:"requestedUser"`
	Title                            string    `json:"title" bson:"title"`
	RequestedDate                    string    `json:"requestedDate" bson:"requestedDate"`
	Type                             string    `json:"type" bson:"type"`
	AdditionalInformation            string    `json:"additionalInformation" bson:"additionalInformation"`
	LongDate                         string    `json:"longDate" bson:"longDate"`
	ShortDate                        string    `json:"shortDate" bson:"shortDate"`
	LongTime                         string    `json:"longTime" bson:"longTime"`
	ShortTime                        string    `json:"shortTime" bson:"shortTime"`
	Overview                         string    `json:"overview" bson:"overview"`
	Year                             string    `json:"year" bson:"year"`
	EpisodesList                     string    `json:"episodesList" bson:"episodesList"`
	SeasonsList                      string    `json:"seasonsList" bson:"seasonsList"`
	PosterImage                      string    `json:"posterImage" bson:"posterImage"`
	ApplicationName                  string    `json:"applicationName" bson:"applicationName"`
	ApplicationURL                   string    `json:"applicationUrl" bson:"applicationUrl"`
	IssueDescription                 string    `json:"issueDescription" bson:"issueDescription"`
	IssueCategory                    string    `json:"issueCategory" bson:"issueCategory"`
	IssueStatus                      string    `json:"issueStatus" bson:"issueStatus"`
	IssueSubject                     string    `json:"issueSubject" bson:"issueSubject"`
	NewIssueComment                  string    `json:"newIssueComment" bson:"newIssueComment"`
	IssueUser                        string    `json:"issueUser" bson:"issueUser"`
	UserName                         string    `json:"userName" bson:"userName"`
	Alias                            string    `json:"alias" bson:"alias"`
	RequestedByAlias                 string    `json:"requestedByAlias" bson:"requestedByAlias"`
	UserPreference                   string    `json:"userPreference" bson:"userPreference"`
	DenyReason                       string    `json:"denyReason" bson:"denyReason"`
	AvailableDate                    string    `json:"availableDate" bson:"availableDate"`
	RequestStatus                    string    `json:"requestStatus" bson:"requestStatus"`
	ProviderID                       string    `json:"providerId" bson:"providerId"`
	PartiallyAvailableEpisodeNumbers string    `json:"partiallyAvailableEpisodeNumbers" bson:"partiallyAvailableEpisodeNumbers"`
	PartiallyAvailableSeasonNumber   string    `json:"partiallyAvailableSeasonNumber" bson:"partiallyAvailableSeasonNumber"`
	PartiallyAvailableEpisodesList   string    `json:"partiallyAvailableEpisodesList" bson:"partiallyAvailableEpisodesList"`
	PartiallyAvailableEpisodeCount   string    `json:"partiallyAvailableEpisodeCount" bson:"partiallyAvailableEpisodeCount"`
	NotificationType                 string    `json:"notificationType" bson:"notificationType"`
	ServiceName                      string    `json:"serviceName" bson:"serviceName"`
	CreatedAt                        time.Time `json:"createdAt" bson:"createdAt"`
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

// ToJSON converts the struct to JSON
func (p *OmbiWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts the JSON to struct
func (p *OmbiWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

// FromHTTPRequest converts the HTTP request to struct
func (p *OmbiWebhookData) FromHTTPRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}
	p.ServiceName = "ombi"
	p.CreatedAt = time.Now()
	return nil
}

// WhoRequestedTvShow returns the list of users who requested a TV show
func WhoRequestedTvShow(title string, season int, episode int) ([]string, error) {
	collection := database.DB.Collection(database.WebhookCollectionName)
	// Define the aggregation pipeline
	pipeline := []bson.M{
		// Match documents where "type" is "TV Show" and "title" is the specified title
		{
			"$match": bson.M{
				"type":  "TV Show",
				"title": title,
			},
		},
		// Split the episodesList and seasonsList into arrays
		{
			"$addFields": bson.M{
				"episodes": bson.M{"$split": []string{"$episodesList", ","}},
				"seasons":  bson.M{"$split": []string{"$seasonsList", ","}},
			},
		},
		// Unwind the episodes and seasons arrays
		{
			"$unwind": bson.M{
				"path":              "$episodes",
				"includeArrayIndex": "episodeIndex",
			},
		},
		{
			"$unwind": bson.M{
				"path":              "$seasons",
				"includeArrayIndex": "seasonIndex",
			},
		},
		// Group by season and episode, collecting unique requested users
		{
			"$group": bson.M{
				"_id": bson.M{
					"season":  "$seasons",
					"episode": "$episodes",
				},
				"users": bson.M{"$addToSet": "$requestedUser"},
			},
		},
		// Match the desired episode and season numbers
		{
			"$match": bson.M{
				"_id.season":  fmt.Sprintf("%d", season),
				"_id.episode": fmt.Sprintf("%d", episode),
			},
		},
		// Project the users field and exclude the _id
		{
			"$project": bson.M{
				"_id":   0,
				"users": 1,
			},
		},
	}

	// Execute the aggregation pipeline
	cursor, err := collection.Aggregate(database.Ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(database.Ctx)

	// Retrieve the result
	var result struct {
		Users []string `bson:"users"`
	}

	if cursor.Next(database.Ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result.Users, nil
}

// RetrieveUsersFromOmbi returns the list of users from Ombi
func RetrieveUsersFromOmbi(service ServiceData) ([]OmbiUser, error) {
	l := logrus.WithField("function", "RetrieveUsersFromOmbi")

	// Validate that the service config has the host and key
	if _, ok := service.Config["host"]; !ok {
		return nil, fmt.Errorf("service config does not have host")
	}
	if _, ok := service.Config["key"]; !ok {
		return nil, fmt.Errorf("service config does not have key")
	}

	// Create a new service rest driver for the media request service.
	serviceRestDriver := servicerestdriver.NewOmbiRestDriver(string(service.ServiceName), service.Config["host"].(string),
		service.Config["key"].(string), l)

	// Reach out to the media request service and get the user's Discord ID.
	usersListURL := fmt.Sprintf("http://%s/ombi/api/v1/Identity/Users", service.Config["host"].(string))
	getUsersListRequest, err := http.NewRequest(http.MethodGet, usersListURL, nil)
	resp, err := serviceRestDriver.Do(getUsersListRequest)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response into a list of users.
	var users []OmbiUser
	err = json.Unmarshal(body, &users)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return users, nil
}

// RetrieveDiscordIDFromOmbi returns the Discord ID of the user who requested the media
func RetrieveDiscordIDFromOmbi(service ServiceData, user OmbiUser) (string, error) {
	l := logrus.WithField("function", "RetrieveDiscordIDFromOmbi")
	// Create a new service rest driver for the media request service.
	serviceRestDriver := servicerestdriver.NewOmbiRestDriver(string(service.ServiceName), service.Config["host"].(string),
		service.Config["key"].(string), l)

	// Reach out to the media request service and get the user's Discord ID.
	userNotificationPrefsURL := fmt.Sprintf("http://%s/ombi/api/v1/Identity/notificationpreferences/%s", service.Config["host"].(string), user.ID)
	getUserNotificationPrefsRequest, err := http.NewRequest(http.MethodGet, userNotificationPrefsURL, nil)
	if err != nil {
		l.WithFields(logrus.Fields{
			"url":          userNotificationPrefsURL,
			"ombiUserName": user.UserName,
			"ombiUserId":   user.ID,
			"error":        err,
		}).Error(err)
		return "", err
	}
	resp, err := serviceRestDriver.Do(getUserNotificationPrefsRequest)
	if err != nil {
		l.WithFields(logrus.Fields{
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
		l.WithFields(logrus.Fields{
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
		l.WithFields(logrus.Fields{
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
