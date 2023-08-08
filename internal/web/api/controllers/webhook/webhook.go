package webhook

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/web/api"
	"time"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// WebhookEntry is the endpoint that handles the inital request for webhooks and routes down to the service-specific func.
func WebhookEntry(w http.ResponseWriter, r *http.Request) {
	l := logrus.WithFields(logrus.Fields{
		"endpoint": r.URL.Path,
		"service":  r.URL.Query().Get("service"),
	})

	l.Info("Webhook received")
	webhookResponse := api.StatusResponse{}
	serviceType := r.URL.Query().Get("service")

	if r.Body == nil {
		api.RenderError("No request body", l, w, r, nil)
		return
	}

	// Store the raw request in the database as UTF-8
	byts, _ := httputil.DumpRequest(r, true)
	filename := fmt.Sprintf("%s_%s.txt", serviceType, time.Now().Format("2006-01-02_15:04:05"))
	models.AddFileToBucket(models.RawRequestWiresBucket, filename, byts, bson.M{"service": serviceType, "event": r.URL.Query().Get("event")})

	// Fire the hook for the given service, or return an error if the service is invalid
	monitoringService := getService(serviceType)
	if monitoringService.monitor == nil {
		api.RenderError("Invalid service", l, w, r, nil)
		return
	}

	// Fire the hook
	err := monitoringService.fireHooks(l, w, r)
	if err != nil {
		api.RenderError(fmt.Sprintf("There was an issue firing the webhook for service %s", serviceType), l, w, r, err)
		return
	}

	// Hooks successfully fired, return response
	webhookResponse.Status = "success"
	webhookResponse.Message = "Webhook fired successfully"
	webhookResponse.Success = true

	render.JSON(w, r, webhookResponse)
}

type ServiceMonitor interface {
	fire(*logrus.Entry, http.ResponseWriter, *http.Request) error
}

// Executes the functions for data collection & storage.
type MonitoringService struct {
	monitor ServiceMonitor
}

// Run the data collection & storage.
func (m MonitoringService) fireHooks(l *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	// Fire webhooks for specific service
	return m.monitor.fire(l, w, r)
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
