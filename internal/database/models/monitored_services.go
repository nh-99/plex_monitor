package models

import (
	"plex_monitor/internal/database"
)

type MonitoredService struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	ApiKey     string `json:"api_key"`
	BaseUrl    string `json:"base_url"`
}

func GetMonitoredService(identifier string) (MonitoredService, error) {
	var monitoredService = MonitoredService{}

	SQL := `SELECT bin_to_uuid(id) as id, identifier, api_key, base_url FROM monitored_services WHERE identifier = ?;`

	err := database.DB.QueryRow(SQL, identifier).Scan(
		&monitoredService.ID,
		&monitoredService.Identifier,
		&monitoredService.ApiKey,
		&monitoredService.BaseUrl,
	)

	if err != nil {
		return monitoredService, err
	}
	return monitoredService, nil
}
