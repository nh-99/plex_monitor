package dashboard

import (
	"html/template"
	"net/http"
	web "plex_monitor/internal/controllers/interface"
	"plex_monitor/internal/controllers/middleware"
	"plex_monitor/internal/database/models"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type viewData struct {
	AppData web.AppData
}

// ViewDashboard will render a UI that displays a summary of the system.
func ViewDashboard(w http.ResponseWriter, r *http.Request) {
	logrus.Info("ViewDashboard route called")
	parsedTemplate, err := template.ParseFiles("./web/html/base.html", "./web/html/_partials/navbar.html", "web/html/_partials/footer.html", "./web/html/dashboard/view.html")
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"err": err,
			},
		).Info("Error ocurred parsing template")
		return
	}
	user := r.Context().Value(middleware.ContextKeyUserID).(models.User)
	err = parsedTemplate.Execute(w, viewData{
		AppData: web.GetAppData(&user),
	})
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"err": err,
			},
		).Info("Error executing template")
		return
	}
}

type activityData struct {
	AppData       web.AppData
	Activity      []models.ActivityStream
	ActivityCount int64
}

// ViewActivity will render a UI that displays recent events.
func ViewActivity(w http.ResponseWriter, r *http.Request) {
	logrus.Info("ViewActivity route called")

	// Set timezone
	loc, err := time.LoadLocation("America/New_York")
	// handle err
	time.Local = loc // -> this is setting the global timezone

	parsedTemplate, err := template.ParseFiles("./web/html/base.html", "./web/html/_partials/navbar.html", "web/html/_partials/footer.html", "./web/html/dashboard/activity.html")
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"err": err,
			},
		).Info("Error ocurred parsing template")
		return
	}
	user := r.Context().Value(middleware.ContextKeyUserID).(models.User)

	// Get number of results to return (limit) from request
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	if err != nil {
		limit = 10
	}

	// Get page number from request
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 1
	}

	// Calculate offset from page & limit
	offset := (page - 1) * limit

	activity, err := models.GetWebhookDataAsActivityStream(offset, limit)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"err":    err,
				"page":   page,
				"limit":  limit,
				"offset": offset,
			},
		).Info("Error ocurred getting activity from DB")
		return
	}
	activityCount, err := models.GetWebhookDataActivityCount()
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"err":    err,
				"page":   page,
				"limit":  limit,
				"offset": offset,
			},
		).Info("Error ocurred getting activity count from DB")
		return
	}

	err = parsedTemplate.Execute(w, activityData{
		AppData:       web.GetAppData(&user),
		Activity:      activity,
		ActivityCount: activityCount,
	})
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"err": err,
			},
		).Info("Error executing template")
		return
	}
}
