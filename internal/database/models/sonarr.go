package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// Episode represents an episode of the TV series.
type Episode struct {
	ID            int    `json:"id" bson:"id"`
	EpisodeNumber int    `json:"episodeNumber" bson:"episodeNumber"`
	SeasonNumber  int    `json:"seasonNumber" bson:"seasonNumber"`
	Title         string `json:"title" bson:"title"`
	AirDate       string `json:"airDate" bson:"airDate"`
	AirDateUTC    string `json:"airDateUtc" bson:"airDateUtc"`
}

// Series represents information about the TV series.
type Series struct {
	ID       int    `json:"id" bson:"id"`
	Title    string `json:"title" bson:"title"`
	Path     string `json:"path" bson:"path"`
	TvdbID   int    `json:"tvdbId" bson:"tvdbId"`
	TvMazeID int    `json:"tvMazeId" bson:"tvMazeId"`
	ImdbID   string `json:"imdbId" bson:"imdbId"`
	Type     string `json:"type" bson:"type"`
}

// QualityRevision represents the quality revision information.
type QualityRevision struct {
	Version  int  `json:"version" bson:"version"`
	Real     int  `json:"real" bson:"real"`
	IsRepack bool `json:"isRepack" bson:"isRepack"`
}

// Quality represents the quality information.
type Quality struct {
	ID         int             `json:"id" bson:"id"`
	Name       string          `json:"name" bson:"name"`
	Source     string          `json:"source" bson:"source"`
	Resolution int             `json:"resolution" bson:"resolution"`
	Revision   QualityRevision `json:"revision" bson:"revision"`
}

// MediaInfo represents the media information.
type MediaInfo struct {
	ContainerFormat                    string  `json:"containerFormat" bson:"containerFormat"`
	VideoFormat                        string  `json:"videoFormat" bson:"videoFormat"`
	VideoCodecID                       string  `json:"videoCodecID" bson:"videoCodecID"`
	VideoProfile                       string  `json:"videoProfile" bson:"videoProfile"`
	VideoCodecLibrary                  string  `json:"videoCodecLibrary" bson:"videoCodecLibrary"`
	VideoBitrate                       int     `json:"videoBitrate" bson:"videoBitrate"`
	VideoBitDepth                      int     `json:"videoBitDepth" bson:"videoBitDepth"`
	VideoMultiViewCount                int     `json:"videoMultiViewCount" bson:"videoMultiViewCount"`
	VideoColourPrimaries               string  `json:"videoColourPrimaries" bson:"videoColourPrimaries"`
	VideoTransferCharacteristics       string  `json:"videoTransferCharacteristics" bson:"videoTransferCharacteristics"`
	VideoHdrFormat                     string  `json:"videoHdrFormat" bson:"videoHdrFormat"`
	VideoHdrFormatCompatibility        string  `json:"videoHdrFormatCompatibility" bson:"videoHdrFormatCompatibility"`
	Width                              int     `json:"width" bson:"width"`
	Height                             int     `json:"height" bson:"height"`
	AudioFormat                        string  `json:"audioFormat" bson:"audioFormat"`
	AudioCodecID                       string  `json:"audioCodecID" bson:"audioCodecID"`
	AudioCodecLibrary                  string  `json:"audioCodecLibrary" bson:"audioCodecLibrary"`
	AudioAdditionalFeatures            string  `json:"audioAdditionalFeatures" bson:"audioAdditionalFeatures"`
	AudioBitrate                       int     `json:"audioBitrate" bson:"audioBitrate"`
	RunTime                            string  `json:"runTime" bson:"runTime"`
	AudioStreamCount                   int     `json:"audioStreamCount" bson:"audioStreamCount"`
	AudioChannelsContainer             int     `json:"audioChannelsContainer" bson:"audioChannelsContainer"`
	AudioChannelsStream                int     `json:"audioChannelsStream" bson:"audioChannelsStream"`
	AudioChannelPositions              string  `json:"audioChannelPositions" bson:"audioChannelPositions"`
	AudioChannelPositionsTextContainer string  `json:"audioChannelPositionsTextContainer" bson:"audioChannelPositionsTextContainer"`
	AudioChannelPositionsTextStream    string  `json:"audioChannelPositionsTextStream" bson:"audioChannelPositionsTextStream"`
	AudioProfile                       string  `json:"audioProfile" bson:"audioProfile"`
	VideoFPS                           float64 `json:"videoFps" bson:"videoFps"`
	AudioLanguages                     string  `json:"audioLanguages" bson:"audioLanguages"`
	Subtitles                          string  `json:"subtitles" bson:"subtitles"`
	ScanType                           string  `json:"scanType" bson:"scanType"`
	SchemaRevision                     int     `json:"schemaRevision" bson:"schemaRevision"`
}

// EpisodeFileEpisode represents an episode of a TV series.
type EpisodeFileEpisode struct {
	SeriesID                 int    `json:"seriesId" bson:"seriesId"`
	TvdbID                   int    `json:"tvdbId" bson:"tvdbId"`
	EpisodeFileID            int    `json:"episodeFileId" bson:"episodeFileId"`
	SeasonNumber             int    `json:"seasonNumber" bson:"seasonNumber"`
	EpisodeNumber            int    `json:"episodeNumber" bson:"episodeNumber"`
	Title                    string `json:"title" bson:"title"`
	AirDate                  string `json:"airDate" bson:"airDate"`
	AirDateUTC               string `json:"airDateUtc" bson:"airDateUtc"`
	Overview                 string `json:"overview" bson:"overview"`
	Monitored                bool   `json:"monitored" bson:"monitored"`
	AbsoluteEpisodeNumber    int    `json:"absoluteEpisodeNumber" bson:"absoluteEpisodeNumber"`
	UnverifiedSceneNumbering bool   `json:"unverifiedSceneNumbering" bson:"unverifiedSceneNumbering"`
	Ratings                  struct {
		Votes int     `json:"votes" bson:"votes"`
		Value float64 `json:"value" bson:"value"`
	} `json:"ratings" bson:"ratings"`
	Images []struct {
		CoverType string `json:"coverType" bson:"coverType"`
		URL       string `json:"url" bson:"url"`
	} `json:"images" bson:"images"`
	LastSearchTime string `json:"lastSearchTime" bson:"lastSearchTime"`
	EpisodeFile    struct {
		IsLoaded bool `json:"isLoaded" bson:"isLoaded"`
	} `json:"episodeFile" bson:"episodeFile"`
	HasFile bool `json:"hasFile" bson:"hasFile"`
	ID      int  `json:"id" bson:"id"`
}

// EpisodeFileEpisodesContainer represents the container for episodes data.
type EpisodeFileEpisodesContainer struct {
	Value    []Episode `json:"value" bson:"value"`
	IsLoaded bool      `json:"isLoaded" bson:"isLoaded"`
}

// EpisodeLanguage represents the language information.
type EpisodeLanguage struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

// EpisodeSeries represents a TV series.
type EpisodeSeries struct {
	Value struct {
		TvdbID            int    `json:"tvdbId" bson:"tvdbId"`
		TvRageID          int    `json:"tvRageId" bson:"tvRageId"`
		TvMazeID          int    `json:"tvMazeId" bson:"tvMazeId"`
		ImdbID            string `json:"imdbId" bson:"imdbId"`
		Title             string `json:"title" bson:"title"`
		CleanTitle        string `json:"cleanTitle" bson:"cleanTitle"`
		SortTitle         string `json:"sortTitle" bson:"sortTitle"`
		Status            string `json:"status" bson:"status"`
		Overview          string `json:"overview" bson:"overview"`
		AirTime           string `json:"airTime" bson:"airTime"`
		Monitored         bool   `json:"monitored" bson:"monitored"`
		QualityProfileID  int    `json:"qualityProfileId" bson:"qualityProfileId"`
		LanguageProfileID int    `json:"languageProfileId" bson:"languageProfileId"`
		SeasonFolder      bool   `json:"seasonFolder" bson:"seasonFolder"`
		LastInfoSync      string `json:"lastInfoSync" bson:"lastInfoSync"`
		Runtime           int    `json:"runtime" bson:"runtime"`
		Images            []struct {
			CoverType string `json:"coverType" bson:"coverType"`
			URL       string `json:"url" bson:"url"`
		} `json:"images" bson:"images"`
		SeriesType        string `json:"seriesType" bson:"seriesType"`
		Network           string `json:"network" bson:"network"`
		UseSceneNumbering bool   `json:"useSceneNumbering" bson:"useSceneNumbering"`
		TitleSlug         string `json:"titleSlug" bson:"titleSlug"`
		Path              string `json:"path" bson:"path"`
		Year              int    `json:"year" bson:"year"`
		Ratings           struct {
			Votes int     `json:"votes" bson:"votes"`
			Value float64 `json:"value" bson:"value"`
		} `json:"ratings" bson:"ratings"`
		Genres []string `json:"genres" bson:"genres"`
		Actors []struct {
			Name      string `json:"name" bson:"name"`
			Character string `json:"character" bson:"character"`
			Images    []struct {
				CoverType string `json:"coverType" bson:"coverType"`
				URL       string `json:"url" bson:"url"`
			} `json:"images" bson:"images"`
		} `json:"actors" bson:"actors"`
		Certification  string `json:"certification" bson:"certification"`
		Added          string `json:"added" bson:"added"`
		FirstAired     string `json:"firstAired" bson:"firstAired"`
		QualityProfile struct {
			Value struct {
				Name           string `json:"name" bson:"name"`
				UpgradeAllowed bool   `json:"upgradeAllowed" bson:"upgradeAllowed"`
				Cutoff         int    `json:"cutoff" bson:"cutoff"`
				Items          []struct {
					Quality struct {
						ID         int    `json:"id" bson:"id"`
						Name       string `json:"name" bson:"name"`
						Source     string `json:"source" bson:"source"`
						Resolution int    `json:"resolution" bson:"resolution"`
					} `json:"quality" bson:"quality"`
					Items   []interface{} `json:"items" bson:"items"`
					Allowed bool          `json:"allowed" bson:"allowed"`
				} `json:"items" bson:"items"`
				ID int `json:"id" bson:"id"`
			} `json:"value" bson:"value"`
			IsLoaded bool `json:"isLoaded" bson:"isLoaded"`
		} `json:"qualityProfile" bson:"qualityProfile"`
		LanguageProfile struct {
			Value struct {
				Name      string `json:"name" bson:"name"`
				Languages []struct {
					Language struct {
						ID   int    `json:"id" bson:"id"`
						Name string `json:"name" bson:"name"`
					} `json:"language" bson:"language"`
					Allowed bool `json:"allowed" bson:"allowed"`
				} `json:"languages" bson:"languages"`
				UpgradeAllowed bool `json:"upgradeAllowed" bson:"upgradeAllowed"`
				Cutoff         struct {
					ID   int    `json:"id" bson:"id"`
					Name string `json:"name" bson:"name"`
				} `json:"cutoff" bson:"cutoff"`
				ID int `json:"id" bson:"id"`
			} `json:"value" bson:"value"`
			IsLoaded bool `json:"isLoaded" bson:"isLoaded"`
		} `json:"languageProfile" bson:"languageProfile"`
		Seasons []struct {
			SeasonNumber int  `json:"seasonNumber" bson:"seasonNumber"`
			Monitored    bool `json:"monitored" bson:"monitored"`
			Images       []struct {
				CoverType string `json:"coverType" bson:"coverType"`
				URL       string `json:"url" bson:"url"`
			} `json:"images" bson:"images"`
		} `json:"seasons" bson:"seasons"`
		Tags []int `json:"tags" bson:"tags"`
		ID   int   `json:"id" bson:"id"`
	} `json:"value" bson:"value"`
	IsLoaded bool `json:"isLoaded" bson:"isLoaded"`
}

// EpisodeFile represents information about the episode file.
type EpisodeFile struct {
	ID            int                           `json:"id" bson:"id"`
	RelativePath  string                        `json:"relativePath" bson:"relativePath"`
	Path          string                        `json:"path" bson:"path"`
	QualityVer    int                           `json:"qualityVersion" bson:"qualityVersion"`
	ReleaseGroup  string                        `json:"releaseGroup" bson:"releaseGroup"`
	SceneName     string                        `json:"sceneName" bson:"sceneName"`
	Size          int                           `json:"size" bson:"size"`
	Quality       interface{}                   `json:"quality" bson:"quality" bson:"quality"`
	MediaInfo     *MediaInfo                    `json:"mediaInfo,omitempty" bson:"mediaInfo,omitempty"`
	Episodes      *EpisodeFileEpisodesContainer `json:"episodes,omitempty" bson:"episodes,omitempty"`
	Language      *EpisodeLanguage              `json:"language,omitempty" bson:"language,omitempty"`
	EpisodeSeries *EpisodeSeries                `json:"series,omitempty" bson:"series,omitempty"`
}

// TvRelease represents information about the grabbed release.
type TvRelease struct {
	Quality        string `json:"quality" bson:"quality"`
	QualityVersion int    `json:"qualityVersion" bson:"qualityVersion"`
	ReleaseGroup   string `json:"releaseGroup" bson:"releaseGroup"`
	ReleaseTitle   string `json:"releaseTitle" bson:"releaseTitle"`
	Indexer        string `json:"indexer" bson:"indexer"`
	Size           int64  `json:"size" bson:"size"`
}

// SonarrWebhookData represents the JSON data structure.
type SonarrWebhookData struct {
	ID                 string         `json:"-" bson:"_id"`
	Series             Series         `json:"series" bson:"series"`
	Episodes           []Episode      `json:"episodes" bson:"episodes"`
	Release            *TvRelease     `json:"release,omitempty" bson:"release,omitempty"`
	EpisodeFile        *EpisodeFile   `json:"episodeFile,omitempty" bson:"episodeFile,omitempty"`
	IsUpgrade          *bool          `json:"isUpgrade,omitempty" bson:"isUpgrade,omitempty"`
	DownloadClient     *string        `json:"downloadClient,omitempty" bson:"downloadClient,omitempty"`
	DownloadClientType *string        `json:"downloadClientType,omitempty" bson:"downloadClientType,omitempty"`
	DownloadID         *string        `json:"downloadId,omitempty" bson:"downloadId,omitempty"`
	DeletedFiles       *[]EpisodeFile `json:"deletedFiles,omitempty" bson:"deletedFiles,omitempty"`
	DeleteReason       *string        `json:"deleteReason,omitempty" bson:"deleteReason,omitempty"`
	EventType          string         `json:"eventType" bson:"eventType"`
	ServiceName        string         `json:"serviceName" bson:"serviceName"`
	CreatedAt          time.Time      `json:"createdAt" bson:"createdAt"`
}

// ToJSON returns the JSON encoding of the struct.
func (p *SonarrWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON decodes the JSON-encoded data and stores the result in the struct.
func (p *SonarrWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

// FromHTTPRequest decodes the JSON-encoded data from the HTTP request and stores the result in the struct.
func (p *SonarrWebhookData) FromHTTPRequest(r *http.Request) error {
	// Parse the request body into the struct
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}
	p.ServiceName = "sonarr"
	p.CreatedAt = time.Now()
	return nil
}
