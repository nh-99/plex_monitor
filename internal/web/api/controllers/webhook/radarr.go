package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
)

// RadarrEventType
const (
	REPOSITORY_RADARR_WEBHOOK        = "radarr"
	RadarrGrab                string = "Grab"
	RadarrDownload            string = "Download"
	RadarrMovieAdded          string = "MovieAdded"
	RadarrApplicationUpdate   string = "ApplicationUpdate"
)

type RadarrWebhookRequest struct {
	Movie struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Year        int    `json:"year"`
		ReleaseDate string `json:"releaseDate"`
		FolderPath  string `json:"folderPath"`
		TMDBID      int    `json:"tmdbId"`
		IMDBID      string `json:"imdbId"`
	} `json:"movie"`
	RemoteMovie struct {
		TMDBID int    `json:"tmdbId"`
		IMDBID string `json:"imdbId"`
		Title  string `json:"title"`
		Year   int    `json:"year"`
	} `json:"remoteMovie"`
	MovieFile struct {
		ID             int    `json:"id"`
		RelativePath   string `json:"relativePath"`
		Path           string `json:"path"`
		Quality        string `json:"quality"`
		QualityVersion int    `json:"qualityVersion"`
		ReleaseGroup   string `json:"releaseGroup"`
		SceneName      string `json:"sceneName"`
		IndexerFlags   string `json:"indexerFlags"`
		Size           int    `json:"size"`
	} `json:"movieFile"`
	IsUpgrade          bool   `json:"isUpgrade"`
	DownloadClient     string `json:"downloadClient"`
	DownloadClientType string `json:"downloadClientType"`
	DownloadID         string `json:"downloadId"`
	Message            string `json:"message"`
	PreviousVersion    string `json:"previousVersion"`
	NewVersion         string `json:"newVersion"`
	EventType          string `json:"eventType"`
}

type RadarrMonitoringService struct{}

// RadarrWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (rms RadarrMonitoringService) fire(w http.ResponseWriter, r *http.Request) {
	radarrWebhookRequest := RadarrWebhookRequest{}
	err := error(nil)

	err = json.NewDecoder(r.Body).Decode(&radarrWebhookRequest)
	if err != nil {
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Collection("radarr_webhook_data").InsertOne(context.TODO(), models.RadarrWebhookData{})
	if err != nil {
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}
}
