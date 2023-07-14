package models

import (
	"encoding/json"
	"net/http"
)

type OmbiWebhookData struct {
	RequestId                        string `bson:"requestId"`
	RequestedUser                    string `bson:"requestedUser"`
	Title                            string `bson:"title"`
	RequestedDate                    string `bson:"requestedDate"`
	Type                             string `bson:"type"`
	AdditionalInformation            string `bson:"additionalInformation"`
	LongDate                         string `bson:"longDate"`
	ShortDate                        string `bson:"shortDate"`
	LongTime                         string `bson:"longTime"`
	ShortTime                        string `bson:"shortTime"`
	Overview                         string `bson:"overview"`
	Year                             string `bson:"year"`
	EpisodesList                     string `bson:"episodesList"`
	SeasonsList                      string `bson:"seasonsList"`
	PosterImage                      string `bson:"posterImage"`
	ApplicationName                  string `bson:"applicationName"`
	ApplicationUrl                   string `bson:"applicationUrl"`
	IssueDescription                 string `bson:"issueDescription"`
	IssueCategory                    string `bson:"issueCategory"`
	IssueStatus                      string `bson:"issueStatus"`
	IssueSubject                     string `bson:"issueSubject"`
	NewIssueComment                  string `bson:"newIssueComment"`
	IssueUser                        string `bson:"issueUser"`
	UserName                         string `bson:"userName"`
	Alias                            string `bson:"alias"`
	RequestedByAlias                 string `bson:"requestedByAlias"`
	UserPreference                   string `bson:"userPreference"`
	DenyReason                       string `bson:"denyReason"`
	AvailableDate                    string `bson:"availableDate"`
	RequestStatus                    string `bson:"requestStatus"`
	ProviderId                       string `bson:"providerId"`
	PartiallyAvailableEpisodeNumbers string `bson:"partiallyAvailableEpisodeNumbers"`
	PartiallyAvailableSeasonNumber   string `bson:"partiallyAvailableSeasonNumber"`
	PartiallyAvailableEpisodesList   string `bson:"partiallyAvailableEpisodesList"`
	PartiallyAvailableEpisodeCount   string `bson:"partiallyAvailableEpisodeCount"`
	NotificationType                 string `bson:"notificationType"`
}

func (p *OmbiWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *OmbiWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *OmbiWebhookData) FromHTTPRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}
	return nil
}
