package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// StatusResponse is a serializer for a generic status response
type StatusResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Success bool          `json:"success"`
	Data    []interface{} `json:"data"`
}

func RenderError(errorMessage string, l *logrus.Entry, w http.ResponseWriter, r *http.Request, err error) StatusResponse {
	// Log the error
	l.WithFields(logrus.Fields{
		"error": err,
	}).Errorf("Encountered error with request: %s", errorMessage)

	response := StatusResponse{}

	// Construct the response
	response.Status = "error"
	response.Message = errorMessage
	w.WriteHeader(http.StatusBadRequest)

	// Return the response
	render.JSON(w, r, response)
	return response
}
