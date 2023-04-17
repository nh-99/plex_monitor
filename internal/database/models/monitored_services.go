package models

import (
	"database/sql"
	"plex_monitor/internal/database"
)

type MonitoredService struct {
	ID         string         `json:"id"`
	Identifier string         `json:"identifier"`
	ApiKey     sql.NullString `json:"api_key"`
	BaseUrl    string         `json:"base_url"`
	IsScanned  bool           `json:"is_scanned"`
}

func GetMonitoredService(identifier string) (MonitoredService, error) {
	var monitoredService = MonitoredService{}

	SQL := `SELECT bin_to_uuid(id) as id, identifier, api_key, base_url, is_scanned FROM monitored_services WHERE identifier = ?;`

	err := database.DB.QueryRow(SQL, identifier).Scan(
		&monitoredService.ID,
		&monitoredService.Identifier,
		&monitoredService.ApiKey,
		&monitoredService.BaseUrl,
		&monitoredService.IsScanned,
	)

	if err != nil {
		return monitoredService, err
	}
	return monitoredService, nil
}

func GetScannableMonitoredServices() ([]*MonitoredService, error) {
	SQL := `SELECT bin_to_uuid(id) as id, identifier, api_key, base_url, is_scanned FROM monitored_services WHERE is_scanned = 1;`
	rows, err := database.DB.Query(SQL)
	if err != nil {
		return nil, err
	}

	monitoredServices := make([]*MonitoredService, 0)
	for rows.Next() {
		monitoredService := new(MonitoredService)
		err := rows.Scan(
			&monitoredService.ID,
			&monitoredService.Identifier,
			&monitoredService.ApiKey,
			&monitoredService.BaseUrl,
			&monitoredService.IsScanned,
		)
		if err != nil {
			return nil, err
		}
		monitoredServices = append(monitoredServices, monitoredService)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return monitoredServices, nil
}
