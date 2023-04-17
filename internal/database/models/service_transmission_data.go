package models

import (
	"database/sql"
	"plex_monitor/internal/database"
	"time"
)

type ServiceTransmissionData struct {
	ID             string         `json:"id"`
	AddedDate      time.Time      `json:"added_date"`
	Error          int            `json:"error"`
	ErrorString    sql.NullString `json:"error_string"`
	IsFinished     bool           `json:"is_finished"`
	Name           sql.NullString `json:"name"`
	PeersConnected int            `json:"peers_connected"`
	Size           int64          `json:"size"`
}

func GetServiceTransmissionDataByName(name string) (ServiceTransmissionData, error) {
	// Get a user from their email
	var serviceTransmissionData = ServiceTransmissionData{}
	SQL := `SELECT bin_to_uuid(id) as id, added_date, error, error_string, is_finished, name, peers_connected, size, updated_at FROM service_transmission_data WHERE name = ?;`

	err := database.DB.QueryRow(SQL, name).Scan(
		&serviceTransmissionData.ID,
	)

	if err != nil {
		return serviceTransmissionData, err
	}
	return serviceTransmissionData, nil
}

func (s ServiceTransmissionData) Commit() error {
	var err error

	SQL := `
	REPLACE INTO
		service_transmission_data (
			id,
			added_date,
			error,
			error_string,
			is_finished,
			name,
			peers_connected,
			size,
			updated_at
		)
	VALUES (uuid_to_bin(?), ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err = database.DB.Exec(SQL,
		s.ID,
		s.AddedDate,
		s.Error,
		s.ErrorString,
		s.IsFinished,
		s.Name,
		s.PeersConnected,
		s.Size,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
