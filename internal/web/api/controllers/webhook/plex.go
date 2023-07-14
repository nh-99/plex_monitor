package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"

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
func (pms PlexMonitoringService) fire(w http.ResponseWriter, r *http.Request) {
	plexWebhookRequest := PlexMonitoringService{}
	err := error(nil)

	err = json.NewDecoder(r.Body).Decode(&plexWebhookRequest)
	if err != nil {
		logrus.Infof("Invalid JSON data: %s", err.Error())
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Collection("plex_webhook_data").InsertOne(context.TODO(), models.PlexWebhookData{})
	if err != nil {
		logrus.Infof("Unable to write to collection: %s", err.Error())
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}
}
