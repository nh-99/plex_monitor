package webhook

import (
	"fmt"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"

	"github.com/sirupsen/logrus"
)

const (
	// RepositoryOmbiWebhook is the name of the repository for the Ombi webhook
	RepositoryOmbiWebhook = "ombi"
)

// OmbiMonitoringService is the struct for the Ombi webhook
type OmbiMonitoringService struct{}

func (rms OmbiMonitoringService) fire(l *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	l.Info("Firing webhook for Ombi")

	ombiWebhookData := models.OmbiWebhookData{}
	err := ombiWebhookData.FromHTTPRequest(r)
	if err != nil {
		return fmt.Errorf("unable to parse request (bad request data): %s", err)
	}

	_, err = database.DB.Collection(database.WebhookCollectionName).InsertOne(database.Ctx, ombiWebhookData)
	if err != nil {
		return fmt.Errorf("unable to write to database: %s", err)
	}

	return nil
}
