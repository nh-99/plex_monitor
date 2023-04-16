package models

import (
	"plex_monitor/internal/database"
	"time"
)

type ServiceSonarrData struct {
	ID                    int       `json:"id"`
	Title                 string    `json:"title"`
	SeasonNumber          int       `json:"season_number"`
	EpisodeNumber         int       `json:"episode_number"`
	AirDateRaw            string    `json:"-"`
	AirDate               time.Time `json:"air_date"`
	TrackedDownloadStatus string    `json:"tracked_download_status"`
	Status                string    `json:"status"`
	HasFile               bool      `json:"has_file"`
	SeriesAirtime         string    `json:"series_airtime"`
	SeriesOverview        string    `json:"series_overview"`
	SeriesCount           int       `json:"series_count"`
	SeriesTitle           string    `json:"series_title"`
	Repository            string    `json:"repository"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (d ServiceSonarrData) Commit() error {
	var err error

	SQL := `
	REPLACE INTO
		service_sonarr_data (
			id,
			title,
			season_number,
			episode_number,
			air_date,
			tracked_download_status,
			status,
			has_file,
			series_airtime,
			series_overview,
			series_count,
			series_title,
			repository,
			updated_at
		)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err = database.DB.Exec(SQL,
		d.ID,
		d.Title,
		d.SeasonNumber,
		d.EpisodeNumber,
		d.AirDateRaw,
		d.TrackedDownloadStatus,
		d.Status,
		d.HasFile,
		d.SeriesAirtime,
		d.SeriesOverview,
		d.SeriesCount,
		d.SeriesTitle,
		d.Repository,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
