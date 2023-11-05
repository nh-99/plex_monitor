package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plex_monitor/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
