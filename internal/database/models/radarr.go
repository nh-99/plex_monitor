package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// RadarrWebhookData is the struct that represents the data that is sent from Radarr.
type RadarrWebhookData struct {
	ID                 string        `json:"-" bson:"_id"`
	Movie              Movie         `json:"movie" bson:"movie"`
	RemoteMovie        *RemoteMovie  `json:"remoteMovie,omitempty" bson:"remoteMovie,omitempty"`
	Release            *MovieRelease `json:"release,omitempty" bson:"release,omitempty"`
	MovieFile          *MovieFile    `json:"movieFile,omitempty" bson:"movieFile,omitempty"`
	DownloadClient     *string       `json:"downloadClient,omitempty" bson:"downloadClient,omitempty"`
	DownloadClientType *string       `json:"downloadClientType,omitempty" bson:"downloadClientType,omitempty"`
	DownloadID         *string       `json:"downloadId,omitempty" bson:"downloadId,omitempty"`
	CustomFormatInfo   *CustomFormat `json:"customFormatInfo,omitempty" bson:"customFormatInfo,omitempty"`
	EventType          string        `json:"eventType" bson:"eventType"`
	InstanceName       string        `json:"instanceName" bson:"instanceName"`
	ApplicationURL     string        `json:"applicationUrl" bson:"applicationUrl"`
	ServiceName        string        `json:"serviceName" bson:"serviceName"`
	CreatedAt          time.Time     `json:"createdAt" bson:"createdAt"`
}

// Movie is the struct that represents the movie data that is sent from Radarr.
type Movie struct {
	ID          int    `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Year        int    `json:"year" bson:"year"`
	ReleaseDate string `json:"releaseDate" bson:"releaseDate"`
	FolderPath  string `json:"folderPath" bson:"folderPath"`
	TmdbID      int    `json:"tmdbId" bson:"tmdbId"`
	ImdbID      string `json:"imdbId" bson:"imdbId"`
	Overview    string `json:"overview" bson:"overview"`
}

// RemoteMovie is the struct that represents the remote movie data that is sent from Radarr.
type RemoteMovie struct {
	TmdbID int    `json:"tmdbId" bson:"tmdbId"`
	ImdbID string `json:"imdbId" bson:"imdbId"`
	Title  string `json:"title" bson:"title"`
	Year   int    `json:"year" bson:"year"`
}

// MovieFile is the struct that represents the movie file data that is sent from Radarr.
type MovieFile struct {
	ID             int            `json:"id" bson:"id"`
	RelativePath   string         `json:"relativePath" bson:"relativePath"`
	Path           string         `json:"path" bson:"path"`
	Quality        string         `json:"quality" bson:"quality"`
	QualityVersion int            `json:"qualityVersion" bson:"qualityVersion"`
	ReleaseGroup   string         `json:"releaseGroup" bson:"releaseGroup"`
	SceneName      string         `json:"sceneName" bson:"sceneName"`
	IndexerFlags   string         `json:"indexerFlags" bson:"indexerFlags"`
	Size           int64          `json:"size" bson:"size"`
	DateAdded      string         `json:"dateAdded" bson:"dateAdded"`
	MediaInfo      MovieMediaInfo `json:"mediaInfo" bson:"mediaInfo"`
}

// MovieMediaInfo is the struct that represents the media info data that is sent from Radarr.
type MovieMediaInfo struct {
	AudioChannels         float64  `json:"audioChannels" bson:"audioChannels"`
	AudioCodec            string   `json:"audioCodec" bson:"audioCodec"`
	AudioLanguages        []string `json:"audioLanguages" bson:"audioLanguages"`
	Height                int      `json:"height" bson:"height"`
	Width                 int      `json:"width" bson:"width"`
	Subtitles             []string `json:"subtitles" bson:"subtitles"`
	VideoCodec            string   `json:"videoCodec" bson:"videoCodec"`
	VideoDynamicRange     string   `json:"videoDynamicRange" bson:"videoDynamicRange"`
	VideoDynamicRangeType string   `json:"videoDynamicRangeType" bson:"videoDynamicRangeType"`
}

// MovieRelease is the struct that represents the movie release data that is sent from Radarr.
type MovieRelease struct {
	Quality           *string  `json:"quality,omitempty" bson:"quality,omitempty"`
	QualityVersion    *int     `json:"qualityVersion,omitempty" bson:"qualityVersion,omitempty"`
	ReleaseGroup      *string  `json:"releaseGroup,omitempty" bson:"releaseGroup,omitempty"`
	ReleaseTitle      string   `json:"releaseTitle" bson:"releaseTitle"`
	Indexer           string   `json:"indexer" bson:"indexer"`
	Size              int      `json:"size" bson:"size"`
	CustomFormatScore *int     `json:"customFormatScore,omitempty" bson:"customFormatScore,omitempty"`
	CustomFormats     []string `json:"customFormats,omitempty" bson:"customFormats,omitempty"`
	IndexerFlags      []string `json:"indexerFlags,omitempty" bson:"indexerFlags,omitempty"`
}

// CustomFormat is the struct that represents the custom format data that is sent from Radarr.
type CustomFormat struct {
	CustomFormats     []string `json:"customFormats" bson:"customFormats"`
	CustomFormatScore int      `json:"customFormatScore" bson:"customFormatScore"`
}

// ToJSON converts the struct to JSON.
func (p *RadarrWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts the JSON to struct.
func (p *RadarrWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

// FromHTTPRequest converts an HTTP request to the struct.
func (p *RadarrWebhookData) FromHTTPRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}
	p.ServiceName = "radarr"
	p.CreatedAt = time.Now()
	return nil
}
