package datacollector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"plex_monitor/internal/database/models"
	"time"
)

const (
	REPOSITORY_SONARR_QUEUE    = "sonarrQueue"
	REPOSITORY_SONARR_CALENDAR = "sonarrCalendar"
)

type sonarrTvSeries struct {
	Title    string `json:"title"`
	Count    int    `json:"seasonCount"`
	Overview string `json:"overview"`
	AirTime  string `json:"airTime"`
}

type sonarrTvShow struct {
	RawAirDate            string         `json:"airDate"`
	AirDate               time.Time      `json:",omitempty"`
	EpisodeNumber         int            `json:"episodeNumber"`
	Id                    int            `json:"id"`
	SeasonNumber          int            `json:"seasonNumber"`
	Title                 string         `json:"title"`
	Series                sonarrTvSeries `json:"series"`
	HasFile               bool           `json:"hasFile"`
	Status                string         `json:"status"`
	TrackedDownloadStatus string         `json:"trackedDownloadStatus"`
}

type SonarrCalendar struct {
	tvShows []sonarrTvShow
}

type SonarrQueue struct {
	tvShows []sonarrTvShow
}

// Collect the data from Sonarrs calendar endpoint.
func (s *SonarrCalendar) collect() error {
	serviceConfig, err := models.GetMonitoredService(REPOSITORY_SONARR_CALENDAR)
	if err != nil {
		return err
	}

	sonarrUrl := serviceConfig.BaseUrl
	resource := "/calendar/"

	// Calculate times for calendar endpoint date filter
	currentTime := time.Now()
	startTime := currentTime
	weeksToAdd := 4
	endTime := currentTime.Add(time.Hour * 24 * 7 * time.Duration(weeksToAdd))

	// Add params to URL
	params := url.Values{}
	params.Add("start", startTime.Format(time.RFC3339Nano))
	params.Add("end", endTime.Format(time.RFC3339Nano))

	u, _ := url.ParseRequestURI(sonarrUrl)
	u.Path = "/api" + resource
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	// Create HTTP client
	httpClient := http.Client{Timeout: time.Duration(10) * time.Second}
	// Start to setup the request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	// Add API key header
	req.Header.Add("X-Api-Key", serviceConfig.ApiKey.String)

	// Do the HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	// Close the request
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal the JSON response into the tvShows variable on the dependency
	err = json.Unmarshal(body, &s.tvShows)
	if err != nil {
		return err
	}

	// Print out shows
	for _, show := range s.tvShows {
		show.AirDate, _ = time.Parse("2006-01-02 15:04", show.RawAirDate+" "+show.Series.AirTime)
		fmt.Printf("[plex-monitor][%s] %s - %s (S%dE%d) airing on %s\n",
			REPOSITORY_SONARR_CALENDAR,
			show.Series.Title,
			show.Title,
			show.SeasonNumber,
			show.EpisodeNumber,
			show.AirDate.Format(time.RFC3339Nano),
		)
	}

	// No errors 🤠
	return nil
}

// Store the Sonarr calendar data.
func (s *SonarrCalendar) store() error {
	for _, show := range s.tvShows {
		serviceSonarrData := models.ServiceSonarrData{
			ID:                    show.Id,
			Title:                 show.Title,
			SeasonNumber:          show.SeasonNumber,
			EpisodeNumber:         show.EpisodeNumber,
			AirDateRaw:            show.RawAirDate,
			TrackedDownloadStatus: show.TrackedDownloadStatus,
			Status:                show.Status,
			HasFile:               show.HasFile,
			SeriesAirtime:         show.Series.AirTime,
			SeriesOverview:        show.Series.Overview,
			SeriesCount:           show.Series.Count,
			SeriesTitle:           show.Series.Title,
			Repository:            REPOSITORY_SONARR_CALENDAR,
			UpdatedAt:             time.Now(),
		}
		err := serviceSonarrData.Commit()

		if err != nil {
			return err
		}
	}

	return nil
}

// Collect the data from Sonarr's queue endpoint.
func (s *SonarrQueue) collect() error {
	serviceConfig, err := models.GetMonitoredService(REPOSITORY_SONARR_QUEUE)
	if err != nil {
		return err
	}

	sonarrUrl := serviceConfig.BaseUrl
	resource := "/queue/"

	u, _ := url.ParseRequestURI(sonarrUrl)
	u.Path = "/api" + resource
	urlStr := fmt.Sprintf("%v", u)

	// Create HTTP client
	httpClient := http.Client{Timeout: time.Duration(10) * time.Second}
	// Start to setup the request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	// Add API key header
	req.Header.Add("X-Api-Key", serviceConfig.ApiKey.String)

	// Do the HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	// Close the request
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal the JSON response into the tvShows variable on the dependency
	err = json.Unmarshal(body, &s.tvShows)
	if err != nil {
		return err
	}

	// Print out shows
	for _, show := range s.tvShows {
		show.AirDate, _ = time.Parse("2006-01-02 15:04", show.RawAirDate+" "+show.Series.AirTime)
		fmt.Printf("[plex-monitor][%s] %s - %s (%s)\n",
			REPOSITORY_SONARR_QUEUE,
			show.Series.Title,
			show.Status,
			show.TrackedDownloadStatus,
		)
	}

	// No errors 🤠
	return nil
}

// Store Sonarr's queue data in the database.
func (s *SonarrQueue) store() error {
	for _, show := range s.tvShows {
		serviceSonarrData := models.ServiceSonarrData{
			ID:                    show.Id,
			Title:                 show.Title,
			SeasonNumber:          show.SeasonNumber,
			EpisodeNumber:         show.EpisodeNumber,
			AirDateRaw:            "1970-01-01 00:00:00",
			TrackedDownloadStatus: show.TrackedDownloadStatus,
			Status:                show.Status,
			HasFile:               show.HasFile,
			SeriesAirtime:         show.Series.AirTime,
			SeriesOverview:        show.Series.Overview,
			SeriesCount:           show.Series.Count,
			SeriesTitle:           show.Series.Title,
			Repository:            REPOSITORY_SONARR_QUEUE,
			UpdatedAt:             time.Now(),
		}
		err := serviceSonarrData.Commit()

		if err != nil {
			return err
		}
	}

	return nil
}
