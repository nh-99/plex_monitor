package user

import (
	"fmt"
	"net/http"
	web "plex_monitor/internal/controllers/interface"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/fun/inspiration"
	"text/template"

	"github.com/sirupsen/logrus"
)

type loginView struct {
	InspirationalQuote string
	Error              string
	HasError           bool
	AppData            web.AppData
}

// ViewLogin will render a UI for the user to login.
func ViewLogin(w http.ResponseWriter, r *http.Request) {
	logrus.Info("ViewLogin route called")
	parsedTemplate, err := template.ParseFiles("./web/html/base.html", "./web/html/user/login.html")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Error ocurred parsing template")
		return
	}
	view := loginView{
		InspirationalQuote: inspiration.GetInspirationalQuote(),
		AppData:            web.GetAppData(nil),
	}
	err = parsedTemplate.Execute(w, view)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Error executing template")
		return
	}
}

type loginData struct {
	InspirationalQuote string
}

// PerformLogin will render a UI for the user to login.
func PerformLogin(w http.ResponseWriter, r *http.Request) {
	logrus.Info("PerformLogin route called")

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	logrus.WithFields(logrus.Fields{
		"email": email,
	}).Info("PerformLogin route called")

	// Check if the user exists
	user, err := models.GetUser("", email)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Error getting user")
		handleUnauthorized(w)
		return
	}

	// Check if the password is correct
	if !user.CheckPassword(password) {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Invalid password supplied")
		handleUnauthorized(w)
		return
	}

	// Set jwt cookie
	tokenString, err := user.GetBearerToken()
	if err != nil {
		logrus.WithError(err).Error("Could not get the bearer token")
	}
	w.Header().Set("Set-Cookie", fmt.Sprintf("jwt=%s; HttpOnly; SameSite=Strict; Path=/;", tokenString))

	// Redirect to the dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func handleUnauthorized(w http.ResponseWriter) {
	// Set the error message
	w.WriteHeader(http.StatusUnauthorized)

	// Render the login page again
	parsedTemplate, err := template.ParseFiles("./web/html/base.html", "./web/html/user/login.html")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Error ocurred parsing template")
		return
	}
	view := loginView{
		InspirationalQuote: inspiration.GetInspirationalQuote(),
		Error:              "Invalid email or password",
		HasError:           true,
	}
	err = parsedTemplate.Execute(w, view)
}
