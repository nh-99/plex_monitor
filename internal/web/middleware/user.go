package middleware

import (
	"plex_monitor/internal/database/models"
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"
)

type key string

const (
	// ContextKeyUserID is the key used to the the user struct in the HTTP context
	ContextKeyUserID key = "User"

	// ClaimsUserIDKey is the key used to set the user ID in the claims
	ClaimsUserIDKey = "user_id"
)

// CreateUserContext is a middleware that adds the User to the request context. If there is no user logged in,
// an anonymous user will be returned.
func CreateUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userForCtx models.User

		_, claims, _ := jwtauth.FromContext(r.Context())
		if userID, ok := claims[ClaimsUserIDKey]; ok {
			userForCtx, _ = models.GetUser("", userID.(string)) // Convert user ID claim to string
		}

		// Return an anonymous user if the user hasn't been found
		if userForCtx.ID == "" {
			userForCtx = models.GetAnonymousUser()
		}

		userCtx := context.WithValue(r.Context(), ContextKeyUserID, userForCtx)
		next.ServeHTTP(w, r.WithContext(userCtx))

		return
	})
}
