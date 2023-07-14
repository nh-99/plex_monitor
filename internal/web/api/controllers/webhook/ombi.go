package webhook

import (
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/web/api"

	"github.com/sirupsen/logrus"
)

const (
	REPOSITORY_OMBI_WEBHOOK = "ombi"
)

type OmbiWebhookRequest struct {
	RequestId                        string `json:"requestId"`
	RequestedUser                    string `json:"requestedUser"`
	Title                            string `json:"title"`
	RequestedDate                    string `json:"requestedDate"`
	Type                             string `json:"type"`
	AdditionalInformation            string `json:"additionalInformation"`
	LongDate                         string `json:"longDate"`
	ShortDate                        string `json:"shortDate"`
	LongTime                         string `json:"longTime"`
	ShortTime                        string `json:"shortTime"`
	Overview                         string `json:"overview"`
	Year                             string `json:"year"`
	EpisodesList                     string `json:"episodesList"`
	SeasonsList                      string `json:"seasonsList"`
	PosterImage                      string `json:"posterImage"`
	ApplicationName                  string `json:"applicationName"`
	ApplicationUrl                   string `json:"applicationUrl"`
	IssueDescription                 string `json:"issueDescription"`
	IssueCategory                    string `json:"issueCategory"`
	IssueStatus                      string `json:"issueStatus"`
	IssueSubject                     string `json:"issueSubject"`
	NewIssueComment                  string `json:"newIssueComment"`
	IssueUser                        string `json:"issueUser"`
	UserName                         string `json:"userName"`
	Alias                            string `json:"alias"`
	RequestedByAlias                 string `json:"requestedByAlias"`
	UserPreference                   string `json:"userPreference"`
	DenyReason                       string `json:"denyReason"`
	AvailableDate                    string `json:"availableDate"`
	RequestStatus                    string `json:"requestStatus"`
	ProviderId                       string `json:"providerId"`
	PartiallyAvailableEpisodeNumbers string `json:"partiallyAvailableEpisodeNumbers"`
	PartiallyAvailableSeasonNumber   string `json:"partiallyAvailableSeasonNumber"`
	PartiallyAvailableEpisodesList   string `json:"partiallyAvailableEpisodesList"`
	PartiallyAvailableEpisodeCount   string `json:"partiallyAvailableEpisodeCount"`
	NotificationType                 string `json:"notificationType"`
}

type OmbiMonitoringService struct{}

// RadarrWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (rms OmbiMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) {
	l.Info("Firing webhook for Sonarr")

	ombiWebhookData := models.OmbiWebhookData{}
	err := ombiWebhookData.FromHTTPRequest(r)
	if err != nil {
		api.RenderError("Could not parse data (bad request data)", l, w, r, err)
		return
	}

	_, err = database.DB.Collection("ombi_webhook_data").InsertOne(database.Ctx, ombiWebhookData)
	if err != nil {
		api.RenderError("Could not store data", l, w, r, err)
		return
	}
}
