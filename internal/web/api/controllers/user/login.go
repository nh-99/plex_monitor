package user

import (
	"encoding/json"
	"net/http"
	"os"
	"plex_monitor/internal/database/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
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
	var loginRequest LoginRequest
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)
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
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": user.ID, "exp": jwtauth.ExpireIn(1460 * time.Hour)}) // 1460 hours == two months
	loginResponse := LoginResponse{
		Token: tokenString,
	}
	render.JSON(w, r, loginResponse) // A chi router helper for serializing and returning json
}
