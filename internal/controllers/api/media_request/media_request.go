package mediarequest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plex_monitor/internal/controllers/api"
	"plex_monitor/internal/database/models"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type mediaRequest struct {
	Name string `json:"name"`
}

type mediaRequestResponse struct {
	ID                 string `json:"id"`
	api.StatusResponse `json:",inline"`
}

// PerformMediaRequest is the endpoint that allows a user to login
func PerformMediaRequest(w http.ResponseWriter, r *http.Request) {
	l := logrus.NewEntry(logrus.StandardLogger())
	// Get the user from the request
	requestUser, exists := models.UserFromContext(r.Context())
	if !exists {
		api.RenderError(l, w, r, fmt.Errorf("the user could not be pulled from the request context, but should be there"), http.StatusInternalServerError)
		return
	}

	// Add user fields to logger
	l = l.WithFields(logrus.Fields{
		"requestUserID": requestUser.ID,
	})

	// Validate & decode the request
	mediaRequestPayload := mediaRequest{}
	err := json.NewDecoder(r.Body).Decode(&mediaRequestPayload)
	if err != nil {
		api.RenderError(l, w, r, err, http.StatusBadRequest)
		return
	}

	// Add request fields to logger
	l = l.WithFields(logrus.Fields{
		"requestedMediaName": mediaRequestPayload.Name,
	})

	// Create a new media request in the DB
	mediaRequest := models.MediaRequest{
		Name:          mediaRequestPayload.Name,
		CurrentStatus: models.MediaRequestStatusRequested, // We always start at the requested state
		RequestedBy:   requestUser.ID,
	}

	// Save the request
	err = mediaRequest.Save()
	if err != nil {
		api.RenderError(l, w, r, err, http.StatusInternalServerError)
		return
	}

	// Return a basic 200 status response
	statusReponse := mediaRequestResponse{
		ID: mediaRequest.ID.Hex(),
		StatusResponse: api.StatusResponse{
			Status:  "success",
			Message: "The request has been added to the queue",
			Success: true,
		},
	}
	render.JSON(w, r, statusReponse)
}
