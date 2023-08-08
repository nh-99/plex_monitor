package webhook

import (
	"fmt"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"

	"github.com/sirupsen/logrus"
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
func (rms RadarrMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	l.Info("Firing webhook for Radarr")

	radarrWebhookData := models.RadarrWebhookData{}
	err := radarrWebhookData.FromHTTPRequest(r)
	if err != nil {
		return fmt.Errorf("could not parse data (bad request data): %w", err)
	}

	_, err = database.DB.Collection("radarr_webhook_data").InsertOne(database.Ctx, radarrWebhookData)
	if err != nil {
		return fmt.Errorf("could not store data: %w", err)
	}

	return nil
}
