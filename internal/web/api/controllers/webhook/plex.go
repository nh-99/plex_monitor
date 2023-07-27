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

type PlexMonitoringService struct{}

// PlexWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (pms PlexMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) {
	l.Info("Firing webhook for Plex")

	err := r.ParseMultipartForm(128 << 20) // Max size 128MB
	if err != nil {
		api.RenderError("Could not parse multipart form", l, w, r, err)
		return
	}

	plexWebhookRequest := models.PlexWebhookData{}
	err = plexWebhookRequest.FromHTTPRequest(r)
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
