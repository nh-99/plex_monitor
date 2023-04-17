package datacollector

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"plex_monitor/internal/database/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	REPOSITORY_TRANSMISSION = "transmission"
)

type torrent struct {
	AddedDate      int64  `json:"addedDate"`
	Error          int    `json:"error"`
	ErrorString    string `json:"errorString"`
	IsFinished     bool   `json:"isFinished"`
	Name           string `json:"name"`
	PeersConnected int    `json:"peersConnected"`
	Size           int64  `json:"size"`
}

type transmissionJsonResponseArguments struct {
	Torrents []torrent `json:"torrents"`
}

type transmissionJsonResponse struct {
	Arguments transmissionJsonResponseArguments `json:"arguments"`
}

type Transmission struct {
	torrentCount int
	torrents     []torrent
}

func (t *Transmission) collect() error {
	serviceConfig, err := models.GetMonitoredService(REPOSITORY_TRANSMISSION)
	if err != nil {
		panic(err)
	}

	serviceUrl := serviceConfig.BaseUrl

	u, _ := url.ParseRequestURI(serviceUrl)
	u.Path = "/transmission/rpc/"
	urlStr := fmt.Sprintf("%v", u)

	var rpcData = `
    {
    "method":"torrent-get",
      "arguments":{
        "fields":[
        "id",
        "addedDate",
        "name",
        "totalSize",
        "error",
        "errorString",
        "eta",
        "isFinished",
        "isStalled",
        "leftUntilDone",
        "metadataPercentComplete",
        "peersConnected",
        "peersGettingFromUs",
        "peersSendingToUs",
        "percentDone",
        "queuePosition",
        "rateDownload",
        "rateUpload",
        "recheckProgress",
        "seedRatioMode",
        "seedRatioLimit",
        "sizeWhenDone",
        "status",
        "trackers",
        "downloadDir",
        "uploadedEver",
        "uploadRatio",
        "webseedsSendingToUs"
        ]
      }
    }
  `
	// Create HTTP client
	httpClient := http.Client{Timeout: time.Duration(10) * time.Second}
	// Start to setup the request
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(rpcData))
	if err != nil {
		panic(err)
	}

	// JSON data
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Do the HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 409 {
		// Retry the request w/ new transmission header
		req, err = http.NewRequest(http.MethodPost, urlStr, strings.NewReader(rpcData))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Add("X-Transmission-Session-Id", resp.Header["X-Transmission-Session-Id"][0])
		resp, err = httpClient.Do(req)
		if err != nil {
			panic(err)
		}
	}

	// Close the request
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON response
	jsonData := transmissionJsonResponse{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		panic(err)
	}

	// Setup data on dependency for storage
	t.torrentCount = len(jsonData.Arguments.Torrents)
	t.torrents = jsonData.Arguments.Torrents

	// No errors 🤠
	return nil
}

func (t *Transmission) store() error {
	for _, torrent := range t.torrents {
		serviceTransmissionData, _ := models.GetServiceTransmissionDataByName(torrent.Name)

		if serviceTransmissionData.ID == "" {
			serviceTransmissionData.ID = uuid.New().String()
		}

		// Set new values on data
		serviceTransmissionData.AddedDate = time.Unix(torrent.AddedDate, 0)
		serviceTransmissionData.Error = torrent.Error
		serviceTransmissionData.ErrorString.String = torrent.ErrorString
		serviceTransmissionData.IsFinished = torrent.IsFinished
		serviceTransmissionData.Name = sql.NullString{String: torrent.Name, Valid: true}
		serviceTransmissionData.PeersConnected = torrent.PeersConnected
		serviceTransmissionData.Size = torrent.Size

		err := serviceTransmissionData.Commit()

		if err != nil {
			panic(err)
		}
	}

	return nil
}
