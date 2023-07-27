package models

import (
	"encoding/json"
	"net/http"
)

type PlexWebhookData struct {
	Event   string `json:"event" bson:"event"`
	User    bool   `json:"user" bson:"user"`
	Owner   bool   `json:"owner" bson:"owner"`
	Account struct {
		ID    int    `json:"id" bson:"id"`
		Thumb string `json:"thumb" bson:"thumb"`
		Title string `json:"title" bson:"title"`
	} `json:"Account" bson:"Account"`
	Server struct {
		Title string `json:"title" bson:"title"`
		UUID  string `json:"uuid" bson:"uuid"`
	} `json:"Server" bson:"Server"`
	Player struct {
		Local         bool   `json:"local" bson:"local"`
		PublicAddress string `json:"publicAddress" bson:"publicAddress"`
		Title         string `json:"title" bson:"title"`
		UUID          string `json:"uuid" bson:"uuid"`
	} `json:"Player" bson:"Player"`
	Metadata struct {
		LibrarySectionType   string `json:"librarySectionType" bson:"librarySectionType"`
		RatingKey            string `json:"ratingKey" bson:"ratingKey"`
		Key                  string `json:"key" bson:"key"`
		ParentRatingKey      string `json:"parentRatingKey" bson:"parentRatingKey"`
		GrandparentRatingKey string `json:"grandparentRatingKey" bson:"grandparentRatingKey"`
		GUID                 string `json:"guid" bson:"guid"`
		LibrarySectionID     int    `json:"librarySectionID" bson:"librarySectionID"`
		Type                 string `json:"type" bson:"type"`
		Title                string `json:"title" bson:"title"`
		GrandparentKey       string `json:"grandparentKey" bson:"grandparentKey"`
		ParentKey            string `json:"parentKey" bson:"parentKey"`
		GrandparentTitle     string `json:"grandparentTitle" bson:"grandparentTitle"`
		ParentTitle          string `json:"parentTitle" bson:"parentTitle"`
		Summary              string `json:"summary" bson:"summary"`
		Index                int    `json:"index" bson:"index"`
		ParentIndex          int    `json:"parentIndex" bson:"parentIndex"`
		RatingCount          int    `json:"ratingCount" bson:"ratingCount"`
		Thumb                string `json:"thumb" bson:"thumb"`
		Art                  string `json:"art" bson:"art"`
		ParentThumb          string `json:"parentThumb" bson:"parentThumb"`
		GrandparentThumb     string `json:"grandparentThumb" bson:"grandparentThumb"`
		GrandparentArt       string `json:"grandparentArt" bson:"grandparentArt"`
		AddedAt              int    `json:"addedAt" bson:"addedAt"`
		UpdatedAt            int    `json:"updatedAt" bson:"updatedAt"`
	} `json:"Metadata" bson:"Metadata"`
}

func (p *PlexWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PlexWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlexWebhookData) FromHTTPRequest(r *http.Request) error {
	// Get the data from the "payload" form field and unmarshal it into the PlexWebhookData struct
	// This is because Plex sends the data as a form field instead of a JSON body
	jsonString := r.FormValue("payload")
	// Unmarshal the JSON string into the PlexWebhookData struct
	err := json.Unmarshal([]byte(jsonString), p)
	if err != nil {
		return err
	}

	return nil
}
