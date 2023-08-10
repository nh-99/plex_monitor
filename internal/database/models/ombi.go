package models

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	// OmbiCollectionName is the name of the collection in the database for Ombi
	OmbiCollectionName = "ombi_webhook_data"
)

// OmbiWebhookData is the struct that represents the data sent by Ombi
type OmbiWebhookData struct {
	ID                               string    `json:"-" bson:"_id"`
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
	p.CreatedAt = time.Now()
	return nil
}
