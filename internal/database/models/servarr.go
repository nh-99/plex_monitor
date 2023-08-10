package models

import (
	"encoding/json"
	"net/http"
)

// Servarr holds structs that are common to the data that is sent from Radarr, Sonarr, and Lidarr.

// ServarrHealthData is the struct that represents the health data that is sent from Radarr.
type ServarrHealthData struct {
	Level          string  `json:"level" bson:"level"`
	Message        string  `json:"message" bson:"message"`
	Type           string  `json:"type" bson:"type"`
	WikiURL        string  `json:"wikiUrl" bson:"wikiUrl"`
	EventType      string  `json:"eventType" bson:"eventType"`
	InstanceName   *string `json:"instanceName,omitempty" bson:"instanceName,omitempty"`
	ApplicationURL *string `json:"applicationUrl,omitempty" bson:"applicationUrl,omitempty"`
}

// FromHTTPRequest converts an HTTP request to the struct.
func (p *ServarrHealthData) FromHTTPRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}
	return nil
}
