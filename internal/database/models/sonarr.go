package models

import (
	"encoding/json"
	"net/http"
)

type SonarrWebhookData struct {
	Series struct {
		ID       int    `bson:"id"`
		Title    string `bson:"title"`
		Path     string `bson:"path"`
		TVDBID   int    `bson:"tvdbId"`
		TVMazeID int    `bson:"tvMazeId"`
		IMDBID   string `bson:"imdbId"`
		Type     string `bson:"type"`
	} `bson:"series"`
	Episodes []struct {
		ID            int    `bson:"id"`
		EpisodeNumber int    `bson:"episodeNumber"`
		SeasonNumber  int    `bson:"seasonNumber"`
		Title         string `bson:"title"`
		AirDate       string `bson:"airDate"`
		AirDateUtc    string `bson:"airDateUtc"`
	} `bson:"episodes"`
	Release struct {
		Quality        string `bson:"quality"`
		QualityVersion int    `bson:"qualityVersion"`
		ReleaseGroup   string `bson:"releaseGroup"`
		ReleaseTitle   string `bson:"releaseTitle"`
		Indexer        string `bson:"indexer"`
		Size           int    `bson:"size"`
	} `bson:"release"`
	DownloadClient     string `bson:"downloadClient"`
	DownloadClientType string `bson:"downloadClientType"`
	DownloadID         string `bson:"downloadId"`
	EventType          string `bson:"eventType"`
}

func (p *SonarrWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *SonarrWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *SonarrWebhookData) FromHTTPRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}
	return nil
}
