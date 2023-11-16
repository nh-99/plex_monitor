package user

import (
	"encoding/json"
	"net/http"
	"plex_monitor/internal/controllers/api"
	"plex_monitor/internal/database/models"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// LoginRequest is the un-serializer for the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the serializer for the login response
type LoginResponse struct {
	Token string `json:"access_token"`
}

// PerformLogin is the endpoint that allows a user to login
func PerformLogin(w http.ResponseWriter, r *http.Request) {
	l := logrus.NewEntry(logrus.StandardLogger())

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}

	// If the email is empty, return an error
	if loginRequest.Email == "" {
		http.Error(w, "No email specified", http.StatusNotFound)
		return
	}

	user, err := models.GetUser("", loginRequest.Email)
	if err != nil {
		http.Error(w, "Incorrect username or password", http.StatusBadRequest)
		return
	}

	// Check if the user is using the correct password
	correctPassword := user.CheckPassword(loginRequest.Password)
	if !correctPassword {
		http.Error(w, "Incorrect username or password", http.StatusForbidden)
		return
	}

	// Encode a JWT auth token for the user
	tokenString, err := user.GetBearerToken()
	if err != nil {
		api.RenderError(l, w, r, err, http.StatusInternalServerError)
		return
	}
	loginResponse := LoginResponse{
		Token: tokenString,
	}
	render.JSON(w, r, loginResponse) // A chi router helper for serializing and returning json
}
