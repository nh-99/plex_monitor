package webhook

import (
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/web/api"

	"github.com/sirupsen/logrus"
)

const (
	REPOSITORY_PLEX_NAME = "plex"
)

type PlexWebhookRequest struct {
	Event   string `json:"event"`
	User    bool   `json:"user"`
	Owner   bool   `json:"owner"`
	Account struct {
		ID    int    `json:"id"`
		Thumb string `json:"thumb"`
		Title string `json:"title"`
	} `json:"Account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"Server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"publicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"Player"`
	Metadata struct {
		LibrarySectionType   string `json:"librarySectionType"`
		RatingKey            string `json:"ratingKey"`
		Key                  string `json:"key"`
		ParentRatingKey      string `json:"parentRatingKey"`
		GrandparentRatingKey string `json:"grandparentRatingKey"`
		GUID                 string `json:"guid"`
		LibrarySectionID     int    `json:"librarySectionID"`
		Type                 string `json:"type"`
		Title                string `json:"title"`
		GrandparentKey       string `json:"grandparentKey"`
		ParentKey            string `json:"parentKey"`
		GrandparentTitle     string `json:"grandparentTitle"`
		ParentTitle          string `json:"parentTitle"`
		Summary              string `json:"summary"`
		Index                int    `json:"index"`
		ParentIndex          int    `json:"parentIndex"`
		RatingCount          int    `json:"ratingCount"`
		Thumb                string `json:"thumb"`
		Art                  string `json:"art"`
		ParentThumb          string `json:"parentThumb"`
		GrandparentThumb     string `json:"grandparentThumb"`
		GrandparentArt       string `json:"grandparentArt"`
		AddedAt              int    `json:"addedAt"`
		UpdatedAt            int    `json:"updatedAt"`
	} `json:"Metadata"`
}

type PlexMonitoringService struct{}

// PlexWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (pms PlexMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) {
	l.Info("Firing webhook for Plex")

	plexWebhookRequest := models.PlexWebhookData{}
	err := plexWebhookRequest.FromHTTPRequest(r)
	if err != nil {
		api.RenderError("Unable to parse request (bad request data)", l, w, r, err)
		return
	}

	_, err = database.DB.Collection("plex_webhook_data").InsertOne(database.Ctx, plexWebhookRequest)
	if err != nil {
		api.RenderError("Unable to write to database", l, w, r, err)
		return
	}
}
