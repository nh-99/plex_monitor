package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// StatusResponse is a serializer for a generic status response
type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// RenderError is a helper function to render an error response and log the error.
func RenderError(l *logrus.Entry, w http.ResponseWriter, r *http.Request, err error, status int) StatusResponse {
	// Log the error
	l.WithError(err).WithFields(logrus.Fields{
		"status": status,
		"method": r.Method,
		"path":   r.URL.Path,
	}).Error("Encountered error with request")

	response := StatusResponse{}

	// Construct the response
	response.Status = "error"
	response.Message = err.Error()
	response.Success = false
	w.WriteHeader(http.StatusBadRequest)

	// Return the response
	render.JSON(w, r, response)
	return response
}
