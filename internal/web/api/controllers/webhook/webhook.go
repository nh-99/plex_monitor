package webhook

import (
	"bytes"
	"io"
	"net/http"
	"plex_monitor/internal/database"
	"plex_monitor/internal/web/api"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// WebhookResponse is the serializer for the login response
type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// WebhookEntry is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func WebhookEntry(w http.ResponseWriter, r *http.Request) {
	l := logrus.WithFields(logrus.Fields{
		"endpoint": r.URL.Path,
		"service":  r.URL.Query().Get("service"),
	})

	l.Info("Webhook received")
	webhookResponse := WebhookResponse{}
	serviceType := r.URL.Query().Get("service")

	if r.Body == nil {
		api.RenderError("No request body", l, w, r, nil)
		return
	}

	// Get the body data as a string for reuse
	requestData, err := io.ReadAll(r.Body)
	if err != nil {
		api.RenderError("Unable to read request body", l, w, r, err)
		return
	}

	// Store the raw request in the database as UTF-8
	_, err = database.DB.Collection("raw_requests").InsertOne(database.Ctx, bson.M{"data": string(requestData), "service": serviceType})
	if err != nil {
		api.RenderError("Unable to write raw data to database", l, w, r, err)
		return
	}

	// Fire the hook for the given service, or return an error if the service is invalid
	monitoringService := getService(serviceType)
	if monitoringService.monitor == nil {
		api.RenderError("Invalid service", l, w, r, nil)
		return
	}

	// Re-construct the body data so we can re-use it & fire the hook
	r.Body = io.NopCloser(bytes.NewReader(requestData))
	monitoringService.fireHooks(l, w, r)

	// Hooks successfully fired, return response
	webhookResponse.Status = "success"
	webhookResponse.Message = "Webhook fired successfully"

	render.JSON(w, r, webhookResponse)
}

type ServiceMonitor interface {
	fire(*logrus.Entry, http.ResponseWriter, *http.Request)
}

// Executes the functions for data collection & storage.
type MonitoringService struct {
	monitor ServiceMonitor
}

// Run the data collection & storage.
func (m MonitoringService) fireHooks(l *logrus.Entry, w http.ResponseWriter, r *http.Request) {
	// Fire webhooks for specific service
	m.monitor.fire(l, w, r)
}

func getService(svcName string) MonitoringService {
	switch svcName {
	case REPOSITORY_PLEX_NAME:
		return MonitoringService{
			monitor: PlexMonitoringService{},
		}
	case REPOSITORY_RADARR_WEBHOOK:
		return MonitoringService{
			monitor: RadarrMonitoringService{},
		}
	case REPOSITORY_SONARR_WEBHOOK:
		return MonitoringService{
			monitor: SonarrMonitoringService{},
		}
	case REPOSITORY_OMBI_WEBHOOK:
		return MonitoringService{
			monitor: OmbiMonitoringService{},
		}
	default:
		return MonitoringService{}
	}
}
