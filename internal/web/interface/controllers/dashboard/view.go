package dashboard

import (
	"html/template"
	"log"
	"net/http"
)

type Dashboard struct {
	CdnUrl string
}

func ViewDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := Dashboard{
		CdnUrl: "http://localhost:8080/static",
	}
	parsedTemplate, _ := template.ParseFiles("./web/templates/dashboard/view.html")
	err := parsedTemplate.Execute(w, dashboard)
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
