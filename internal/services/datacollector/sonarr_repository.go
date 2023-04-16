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

type sonarrTvSeries struct {
	Title    string `json:"title"`
	Count    int    `json:"seasonCount"`
	Overview string `json:"overview"`
	AirTime  string `json:"airTime"`
}

type sonarrTvShow struct {
	AirDate       string         `json:"air_date"`
	EpisodeNumber int            `json:"episodeNumber"`
	Id            int            `json:"id"`
	SeasonNumber  int            `json:"seasonNumber"`
	Title         string         `json:"title"`
	Series        sonarrTvSeries `json:"series"`
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
		fmt.Printf("%s - %s (S%dE%d)\n", show.Series.Title, show.Title, show.SeasonNumber, show.EpisodeNumber)
	}

	// No errors 🤠
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
