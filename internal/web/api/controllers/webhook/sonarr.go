package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
)

// SonarrEventType
const (
	REPOSITORY_SONARR_WEBHOOK        = "sonarr"
	SonarrGrab                string = "Grab"
	SonarrDownload            string = "Download"
	SonarrMovieAdded          string = "MovieAdded"
	SonarrApplicationUpdate   string = "ApplicationUpdate"
)

type SonarrWebhookRequest struct {
	Series struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Path     string `json:"path"`
		TVDBID   int    `json:"tvdbId"`
		TVMazeID int    `json:"tvMazeId"`
		IMDBID   string `json:"imdbId"`
		Type     string `json:"type"`
	} `json:"series"`
	Episodes []struct {
		ID            int    `json:"id"`
		EpisodeNumber int    `json:"episodeNumber"`
		SeasonNumber  int    `json:"seasonNumber"`
		Title         string `json:"title"`
		AirDate       string `json:"airDate"`
		AirDateUtc    string `json:"airDateUtc"`
	} `json:"episodes"`
	Release struct {
		Quality        string `json:"quality"`
		QualityVersion int    `json:"qualityVersion"`
		ReleaseGroup   string `json:"releaseGroup"`
		ReleaseTitle   string `json:"releaseTitle"`
		Indexer        string `json:"indexer"`
		Size           int    `json:"size"`
	} `json:"release"`
	DownloadClient     string `json:"downloadClient"`
	DownloadClientType string `json:"downloadClientType"`
	DownloadID         string `json:"downloadId"`
	EventType          string `json:"eventType"`
}

type SonarrMonitoringService struct{}

// SonarrWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (rms SonarrMonitoringService) fire(w http.ResponseWriter, r *http.Request) {
	sonarrWebhookRequest := SonarrWebhookRequest{}
	err := error(nil)

	err = json.NewDecoder(r.Body).Decode(&sonarrWebhookRequest)
	if err != nil {
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Collection("sonarr_webhook_data").InsertOne(context.TODO(), models.SonarrWebhookData{})
	if err != nil {
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}
}
