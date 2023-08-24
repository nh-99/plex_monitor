package user

import (
	"net/http"
	"text/template"

	"github.com/sirupsen/logrus"
)

func ViewDashboard(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("./web/templates/dashboard/view.html")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		logrus.Errorf("Error executing template: %s", err.Error())
		return
	}
}
