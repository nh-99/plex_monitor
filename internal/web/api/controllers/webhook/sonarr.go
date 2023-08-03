package webhook

import (
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/web/api"

	"github.com/sirupsen/logrus"
)

// SonarrEventType
const (
	REPOSITORY_SONARR_WEBHOOK        = "sonarr"
	SonarrGrab                string = "Grab"
	SonarrDownload            string = "Download"
	SonarrMovieAdded          string = "MovieAdded"
	SonarrApplicationUpdate   string = "ApplicationUpdate"
)

type SonarrMonitoringService struct{}

// SonarrWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (rms SonarrMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) {
	l.Info("Firing webhook for Sonarr")

	sonarrWebhookData := models.SonarrWebhookData{}
	err := sonarrWebhookData.FromHTTPRequest(r)

	if err != nil {
		api.RenderError("Could not parse request (bad request data)", l, w, r, err)
		return
	}

	_, err = database.DB.Collection("sonarr_webhook_data").InsertOne(database.Ctx, sonarrWebhookData)
	if err != nil {
		api.RenderError("Could not store data", l, w, r, err)
		return
	}
}
