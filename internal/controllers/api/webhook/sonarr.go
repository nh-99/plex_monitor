package webhook

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/pipeline"
	"strings"

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

	// Parse the request body into a string (in case we need to re-process the request)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("could not parse request body: %w", err)
	}

	// Set the request body back to the original so we can parse it again
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	sonarrWebhookData := models.SonarrWebhookData{}
	err = sonarrWebhookData.FromHTTPRequest(r)

	if err != nil {
		return fmt.Errorf("unable to parse request (bad request data): %s", err)
	}

	// If the event type contains "Health", then we need to parse the data differently.
	if strings.Contains(sonarrWebhookData.EventType, "Health") {
		// Set the request body back to the original so we can parse it again
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Parse the data into the health struct
		healthData := models.ServarrHealthData{}
		err = healthData.FromHTTPRequest(r)
		if err != nil {
			return fmt.Errorf("could not parse data (bad request data): %w", err)
		}

		healthData.ServiceName = "sonarr"

		// Store the data in the database
		_, err := database.DB.Collection(database.WebhookCollectionName).InsertOne(database.Ctx, healthData)
		if err != nil {
			return fmt.Errorf("could not store data: %w", err)
		}

		// Return early since we don't need to do anything else
		return nil
	}

	_, err = database.DB.Collection(database.WebhookCollectionName).InsertOne(database.Ctx, sonarrWebhookData)
	if err != nil {
		return fmt.Errorf("unable to write to database: %s", err)
	}

	// Run the pipelines
	err = runSonarrPipelines(sonarrWebhookData, l)
	if err != nil {
		return fmt.Errorf("unable to run pipelines: %w", err)
	}

	return nil
}

func runSonarrPipelines(sonarrWebhookData models.SonarrWebhookData, l *logrus.Entry) error {
	// Run the pipeline step for Radarr
	pipelineID := pipeline.GeneratePipelineID("TV", sonarrWebhookData.Series.Title)
	pipelineData, err := pipeline.GetOrCreateMediaRequestPipeline(pipelineID)
	if err != nil {
		return fmt.Errorf("unable to get pipeline: %w", err)
	}

	switch sonarrWebhookData.EventType {
	case RadarrMovieDownloadEventType:
		go func() {
			pipelineStep := pipeline.MediaRequestImported
			err = pipelineData.RunStep(pipelineStep)
			if err != nil {
				l.WithField("err", err).Errorf("unable to run pipeline step: %v", err)
			}
		}()
	case RadarrMovieAddedEventType:
		go func() {
			pipelineStep := pipeline.MediaRequestIngestedBySonarr
			err := pipelineData.MarkStepAsSkipped(pipeline.MediaRequestIngestedByRadarr)
			if err != nil {
				l.WithField("err", err).Errorf("unable to mark step as skipped: %v", err)
			}
			err = pipelineData.RunStep(pipelineStep)
			if err != nil {
				l.WithField("err", err).Errorf("unable to run pipeline step: %v", err)
			}
		}()
	}

	return nil
}
