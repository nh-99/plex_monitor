package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// PlexWebhookData is the struct that represents the data sent by Plex webhooks
type PlexWebhookData struct {
	ID      string `json:"-" bson:"_id"`
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
		LibrarySectionType    string  `json:"librarySectionType" bson:"librarySectionType"`
		RatingKey             string  `json:"ratingKey" bson:"ratingKey"`
		Key                   string  `json:"key" bson:"key"`
		MetaGUID              string  `json:"guid" bson:"guid"`
		Studio                string  `json:"studio" bson:"studio"`
		Type                  string  `json:"type" bson:"type"`
		Title                 string  `json:"title" bson:"title"`
		LibrarySectionTitle   string  `json:"librarySectionTitle" bson:"librarySectionTitle"`
		LibrarySectionID      int     `json:"librarySectionID" bson:"librarySectionID"`
		LibrarySectionKey     string  `json:"librarySectionKey" bson:"librarySectionKey"`
		ContentRating         string  `json:"contentRating" bson:"contentRating"`
		Summary               string  `json:"summary" bson:"summary"`
		NumericRating         float64 `json:"rating" bson:"rating"`
		AudienceRating        float64 `json:"audienceRating" bson:"audienceRating"`
		ViewOffset            int     `json:"viewOffset" bson:"viewOffset"`
		LastViewedAt          int     `json:"lastViewedAt" bson:"lastViewedAt"`
		Year                  int     `json:"year" bson:"year"`
		Tagline               string  `json:"tagline" bson:"tagline"`
		Thumb                 string  `json:"thumb" bson:"thumb"`
		Art                   string  `json:"art" bson:"art"`
		Duration              int     `json:"duration" bson:"duration"`
		OriginallyAvailableAt string  `json:"originallyAvailableAt" bson:"originallyAvailableAt"`
		AddedAt               int     `json:"addedAt" bson:"addedAt"`
		UpdatedAt             int     `json:"updatedAt" bson:"updatedAt"`
		AudienceRatingImage   string  `json:"audienceRatingImage" bson:"audienceRatingImage"`
		PrimaryExtraKey       string  `json:"primaryExtraKey" bson:"primaryExtraKey"`
		RatingImage           string  `json:"ratingImage" bson:"ratingImage"`
		Genre                 []struct {
			ID     int    `json:"id" bson:"id"`
			Filter string `json:"filter" bson:"filter"`
			Tag    string `json:"tag" bson:"tag"`
			Count  int    `json:"count" bson:"count"`
		} `json:"Genre" bson:"Genre"`
		Country []struct {
			ID     int    `json:"id" bson:"id"`
			Filter string `json:"filter" bson:"filter"`
			Tag    string `json:"tag" bson:"tag"`
			Count  int    `json:"count" bson:"count"`
		} `json:"Country" bson:"Country"`
		GUID []struct {
			ID string `json:"id" bson:"id"`
		} `json:"Guid" bson:"Guid"`
		Rating []struct {
			Image string  `json:"image" bson:"image"`
			Value float64 `json:"value" bson:"value"`
			Type  string  `json:"type" bson:"type"`
			Count int     `json:"count" bson:"count"`
		} `json:"Rating" bson:"Rating"`
		Director []struct {
			ID     int    `json:"id" bson:"id"`
			Filter string `json:"filter" bson:"filter"`
			Tag    string `json:"tag" bson:"tag"`
			TagKey string `json:"tagKey" bson:"tagKey"`
			Count  int    `json:"count" bson:"count"`
			Thumb  string `json:"thumb" bson:"thumb"`
		} `json:"Director" bson:"Director"`
		Writer []struct {
			ID     int    `json:"id" bson:"id"`
			Filter string `json:"filter" bson:"filter"`
			Tag    string `json:"tag" bson:"tag"`
			TagKey string `json:"tagKey" bson:"tagKey"`
			Count  int    `json:"count" bson:"count"`
			Thumb  string `json:"thumb" bson:"thumb"`
		} `json:"Writer" bson:"Writer"`
		Role []struct {
			ID     int    `json:"id" bson:"id"`
			Filter string `json:"filter" bson:"filter"`
			Tag    string `json:"tag" bson:"tag"`
			TagKey string `json:"tagKey" bson:"tagKey"`
			Count  int    `json:"count" bson:"count"`
			Role   string `json:"role" bson:"role"`
			Thumb  string `json:"thumb" bson:"thumb"`
		} `json:"Role" bson:"Role"`
		Producer []struct {
			ID     int    `json:"id" bson:"id"`
			Filter string `json:"filter" bson:"filter"`
			Tag    string `json:"tag" bson:"tag"`
			TagKey string `json:"tagKey" bson:"tagKey"`
			Count  int    `json:"count" bson:"count"`
			Thumb  string `json:"thumb" bson:"thumb"`
		} `json:"Producer" bson:"Producer"`
	} `json:"Metadata" bson:"Metadata"`
	ServiceName string    `json:"serviceName" bson:"serviceName"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

// ToJSON converts the PlexWebhookData struct to a JSON string
func (p *PlexWebhookData) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts a JSON string to a PlexWebhookData struct
func (p *PlexWebhookData) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

// FromHTTPRequest converts a HTTP request to a PlexWebhookData struct
func (p *PlexWebhookData) FromHTTPRequest(r *http.Request) error {
	// Get the data from the "payload" form field and unmarshal it into the PlexWebhookData struct
	// This is because Plex sends the data as a form field instead of a JSON body
	jsonString := r.FormValue("payload")
	// Unmarshal the JSON string into the PlexWebhookData struct
	err := json.Unmarshal([]byte(jsonString), p)
	if err != nil {
		return err
	}

	p.ServiceName = "plex"
	p.CreatedAt = time.Now()

	return nil
}
