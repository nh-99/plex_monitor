package models

import (
	"encoding/json"
	"net/http"
)

type RadarrWebhookData struct {
	Movie struct {
		ID          int    `bson:"id"`
		Title       string `bson:"title"`
		Year        int    `bson:"year"`
		ReleaseDate string `bson:"releaseDate"`
		FolderPath  string `bson:"folderPath"`
		TMDBID      int    `bson:"tmdbId"`
		IMDBID      string `bson:"imdbId"`
	} `bson:"movie"`
	RemoteMovie struct {
		TMDBID int    `bson:"tmdbId"`
		IMDBID string `bson:"imdbId"`
		Title  string `bson:"title"`
		Year   int    `bson:"year"`
	} `bson:"remoteMovie"`
	MovieFile struct {
		ID             int    `bson:"id"`
		RelativePath   string `bson:"relativePath"`
		Path           string `bson:"path"`
		Quality        string `bson:"quality"`
		QualityVersion int    `bson:"qualityVersion"`
		ReleaseGroup   string `bson:"releaseGroup"`
		SceneName      string `bson:"sceneName"`
		IndexerFlags   string `bson:"indexerFlags"`
		Size           int    `bson:"size"`
	} `bson:"movieFile"`
	IsUpgrade          bool   `bson:"isUpgrade"`
	DownloadClient     string `bson:"downloadClient"`
	DownloadClientType string `bson:"downloadClientType"`
	DownloadID         string `bson:"downloadId"`
	Message            string `bson:"message"`
	PreviousVersion    string `bson:"previousVersion"`
	NewVersion         string `bson:"newVersion"`
	EventType          string `bson:"eventType"`
}

func (p *RadarrWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *RadarrWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *RadarrWebhookData) FromHTTPRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}
	return nil
}
