package webhook

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// RepositoryRadarrWebhook is the name of the repository for the Radarr webhook
	RepositoryRadarrWebhook = "radarr"
	// RadarrMovieAddedEventType is the event type for a new movie added
	RadarrMovieAddedEventType = "MovieAdded"
	// RadarrMovieDownloadEventType is the event type for a movie download
	RadarrMovieDownloadEventType = "Download"
)

// RadarrMonitoringService is the struct for the Radarr webhook
type RadarrMonitoringService struct{}

// RadarrWebhook is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func (rms RadarrMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	l.Info("Firing webhook for Radarr")

	// Parse the request body into a string (in case we need to re-process the request)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("could not parse request body: %w", err)
	}

	// Set the request body back to the original so we can parse it again
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	radarrWebhookData := models.RadarrWebhookData{}
	err = radarrWebhookData.FromHTTPRequest(r)
	if err != nil {
		return fmt.Errorf("could not parse data (bad request data): %w", err)
	}

	// If the event type contains "Health", then we need to parse the data differently.
	if strings.Contains(radarrWebhookData.EventType, "Health") {
		// Set the request body back to the original so we can parse it again
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Parse the data into the health struct
		healthData := models.ServarrHealthData{}
		err = healthData.FromHTTPRequest(r)
		if err != nil {
			return fmt.Errorf("could not parse data (bad request data): %w", err)
		}

		healthData.ServiceName = "radarr"

		// Store the data in the database
		_, err := database.DB.Collection(database.WebhookCollectionName).InsertOne(database.Ctx, healthData)
		if err != nil {
			return fmt.Errorf("could not store data: %w", err)
		}

		// Return early since we don't need to do anything else
		return nil
	}

	_, err = database.DB.Collection(database.WebhookCollectionName).InsertOne(database.Ctx, radarrWebhookData)
	if err != nil {
		return fmt.Errorf("could not store data: %w", err)
	}

	return nil
}
