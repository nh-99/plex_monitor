package middleware

import (
	"fmt"
	"net/http"
	"plex_monitor/internal/controllers/api"
	"plex_monitor/internal/database/models"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sirupsen/logrus"
)

const (
	// ClaimsUserIDKey is the key used to set the user ID in the claims
	ClaimsUserIDKey = "user_id"
)

// CreateUserContext is a middleware that adds the User to the request context. If there is no user logged in,
// an anonymous user will be returned.
func CreateUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userForCtx models.User
		var err error

		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			api.RenderError(logrus.WithField("entry", "UserMiddleware"), w, r, fmt.Errorf("could not get the token from the context"), http.StatusInternalServerError)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			api.RenderError(logrus.WithField("entry", "UserMiddleware"), w, r, fmt.Errorf("invalid token"), http.StatusUnauthorized)
			return
		}

		if userID, ok := claims[ClaimsUserIDKey]; ok {
			fmt.Printf("User ID: %v\n", userID)
			userForCtx, err = models.GetUser(userID.(string), "") // Convert user ID claim to string

			if err != nil {
				logrus.WithFields(logrus.Fields{
					"userID": userID,
				}).WithError(err).Error("Could not get the user")
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"claims": claims,
			}).Error("Could not get the user ID from the claims")
		}

		userCtx := userForCtx.NewContext(r.Context())
		r = r.WithContext(userCtx)
		next.ServeHTTP(w, r)
	})
}
