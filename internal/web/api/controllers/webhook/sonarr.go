package webhook

import (
	"fmt"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"

	"github.com/sirupsen/logrus"
)

const (
	// RepositorySonarrWebhook is the name of the repository for the Sonarr webhook
	RepositorySonarrWebhook = "sonarr"
)

// SonarrMonitoringService is the struct for the Sonarr webhook
type SonarrMonitoringService struct{}

// SonarrWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (rms SonarrMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	l.Info("Firing webhook for Sonarr")

	sonarrWebhookData := models.SonarrWebhookData{}
	err := sonarrWebhookData.FromHTTPRequest(r)

	if err != nil {
		return fmt.Errorf("unable to parse request (bad request data): %s", err)
	}

	_, err = database.DB.Collection("sonarr_webhook_data").InsertOne(database.Ctx, sonarrWebhookData)
	if err != nil {
		return fmt.Errorf("unable to write to database: %s", err)
	}

	return nil
}
