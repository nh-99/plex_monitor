package webhook

import (
	"fmt"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"

	"github.com/sirupsen/logrus"
)

const (
	// RepositoryPlexName is the name of the repository for the Plex webhook
	RepositoryPlexName = "plex"
)

// PlexMonitoringService is the struct for the Plex webhook
type PlexMonitoringService struct{}

// PlexWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (pms PlexMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	l.Info("Firing webhook for Plex")

	err := r.ParseMultipartForm(128 << 20) // Max size 128MB
	if err != nil {
		return fmt.Errorf("unable to parse multipart form: %s", err)
	}

	plexWebhookRequest := models.PlexWebhookData{}
	err = plexWebhookRequest.FromHTTPRequest(r)
	if err != nil {
		return fmt.Errorf("unable to parse request (bad request data): %s", err)
	}

	_, err = database.DB.Collection(models.WebhookCollectionName).InsertOne(database.Ctx, plexWebhookRequest)
	if err != nil {
		return fmt.Errorf("unable to write to database: %s", err)
	}

	return nil
}
