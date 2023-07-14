package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"plex_monitor/internal/database"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// WebhookResponse is the serializer for the login response
type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// WebhookEntry is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func WebhookEntry(w http.ResponseWriter, r *http.Request) {
	webhookResponse := WebhookResponse{}
	serviceType := r.URL.Query().Get("service")

	// Get the body data as a string for reuse
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Construct the response
		webhookResponse.Status = "error"
		webhookResponse.Message = "Request body is empty"
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Unable to read request body: %s", err.Error())

		// Return the response
		render.JSON(w, r, webhookResponse)
		return
	}

	// Store raw response data in mongo
	rawResponse := make(map[string]interface{})
	err = json.NewDecoder(bytes.NewReader(requestData)).Decode(&rawResponse)
	if err != nil {
		// Construct the response
		webhookResponse.Status = "error"
		webhookResponse.Message = "Invalid JSON data"
		w.WriteHeader(http.StatusBadRequest)

		// Read the body so we can log it
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		// Log the body
		logrus.Infof("[plex_monitor] Invalid JSON data: %s", b)

		// Return the response
		render.JSON(w, r, webhookResponse)
		return
	}
	rawResponse["service"] = serviceType
	database.DB.Collection("raw_responses").InsertOne(context.Background(), rawResponse)

	// Fire the hook for the given service, or return an error if the service is invalid
	monitoringService := getService(serviceType)
	if monitoringService.monitor == nil {
		webhookResponse.Status = "error"
		webhookResponse.Message = "Invalid service"
		w.WriteHeader(http.StatusBadRequest)
		logrus.Infof("[plex_monitor] Invalid service attempted: %s", serviceType)
		render.JSON(w, r, webhookResponse)
		return
	}

	// Re-construct the body data so we can re-use it & fire the hook
	r.Body = ioutil.NopCloser(bytes.NewReader(requestData))
	monitoringService.fireHooks(w, r)

	// Hooks successfully fired, return response
	webhookResponse.Status = "success"
	webhookResponse.Message = "Webhook fired successfully"

	render.JSON(w, r, webhookResponse)
}

type ServiceMonitor interface {
	fire(http.ResponseWriter, *http.Request)
}

// Executes the functions for data collection & storage.
type MonitoringService struct {
	monitor ServiceMonitor
}

// Run the data collection & storage.
func (m MonitoringService) fireHooks(w http.ResponseWriter, r *http.Request) {
	// Fire webhooks for specific service
	m.monitor.fire(w, r)
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
