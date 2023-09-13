package webhook

import (
	"fmt"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/pipeline"

	"github.com/sirupsen/logrus"
)

const (
	// RepositoryOmbiWebhook is the name of the repository for the Ombi webhook
	RepositoryOmbiWebhook = "ombi"
	// NewRequestNotificationType is the notification type for a new request
	NewRequestNotificationType = "NewRequest"
)

// OmbiMonitoringService is the struct for the Ombi webhook
type OmbiMonitoringService struct{}

func (rms OmbiMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	l.Info("Firing webhook for Ombi")

	ombiWebhookData := models.OmbiWebhookData{}
	err := ombiWebhookData.FromHTTPRequest(r)
	if err != nil {
		return fmt.Errorf("unable to parse request (bad request data): %w", err)
	}

	_, err = database.DB.Collection(database.WebhookCollectionName).InsertOne(database.Ctx, ombiWebhookData)
	if err != nil {
		return fmt.Errorf("unable to write to database: %w", err)
	}

	// Run the pipelines
	err = runOmbiPipelines(ombiWebhookData, l)
	if err != nil {
		return fmt.Errorf("unable to run pipelines: %w", err)
	}

	return nil
}

func runOmbiPipelines(ombiWebhookData models.OmbiWebhookData, l *logrus.Entry) error {
	// Run the pipeline step for Ombi
	pipelineID := pipeline.GeneratePipelineID(ombiWebhookData.Type, ombiWebhookData.Title)
	pipelineData, err := pipeline.GetOrCreateMediaRequestPipeline(pipelineID)
	pipelineData.AddMetadata("ombi", map[string]interface{}{
		"requestID": ombiWebhookData.RequestID,
		"userName":  ombiWebhookData.UserName,
		"mediaType": ombiWebhookData.Type,
	})
	if err != nil {
		return fmt.Errorf("unable to get pipeline: %w", err)
	}
	err = pipelineData.Save()
	if err != nil {
		return fmt.Errorf("unable to save pipeline: %w", err)
	}

	// If the notification type is "NewRequest", then we need to run the pipeline step.
	switch ombiWebhookData.NotificationType {
	case NewRequestNotificationType:
		// Run the pipeline step
		go func() {
			err = pipelineData.RunStep(pipeline.MediaRequestRequested)
			if err != nil {
				l.WithFields(logrus.Fields{
					"step": pipeline.MediaRequestRequested,
					"err":  err,
				}).Errorf("unable to run pipeline step: %v", err)
			}
		}()
	}

	return nil
}
