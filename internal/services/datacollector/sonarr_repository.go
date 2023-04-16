package datacollector

import (
	"fmt"
	"net/http"
	"net/url"
	"plex_monitor/internal/database/models"
	"time"
)

type sonarrTvShow struct {
	absoluteEpisodeNumber      int
	airDate                    string
	airDateUtc                 string
	episodeFile                []string
	episodeFileId              int
	episodeNumber              int
	hasFile                    bool
	id                         int
	lastSearchTime             string
	monitored                  bool
	overview                   string
	sceneAbsoluteEpisodeNumber int
	sceneEpisodeNumber         int
	sceneSeasonNumber          int
	seasonNumber               int
	series                     []string
	seriesId                   int
	title                      string
	unverifiedSceneNumbering   bool
}

type SonarrCalendar struct {
	tvShows []sonarrTvShow
}

type SonarrQueue struct {
	tvShows []sonarrTvShow
}

func (s SonarrCalendar) collect() error {
	serviceConfig, err := models.GetMonitoredService(SONARR_CALENDAR)
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
	req.Header.Add("X-Api-Key", serviceConfig.ApiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	// Process response into sonarrTvShow's
	// TODO:

	return nil
}

func (s SonarrCalendar) store(db Database) error {
	fmt.Println("calendar store")
	return nil
}

func (s SonarrQueue) collect() error {
	// Get queue
	fmt.Println("queue collect")
	return nil
}

func (s SonarrQueue) store(db Database) error {
	// Store queue
	fmt.Println("queue store")
	return nil
}
