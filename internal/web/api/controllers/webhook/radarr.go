package webhook

import (
	"fmt"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"

	"github.com/sirupsen/logrus"
)

const (
	// RepositoryRadarrWebhook is the name of the repository for the Radarr webhook
	RepositoryRadarrWebhook = "radarr"
)

// RadarrMonitoringService is the struct for the Radarr webhook
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
