package webhook

import (
	"net/http"
	"plex_monitor/internal/web/middleware"

	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()

	// Protected endpoints
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(basicAuth)

		// Custom middleware to add user to request context, for easy access
		r.Use(middleware.CreateUserContext)

		r.Post("/", WebhookEntry)
	})

	return router
}

// basicAuth is a middleware function that checks for valid basic authentication credentials
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the Authorization header is present and starts with "Basic "
		if authHeader != "" && len(authHeader) > 6 && authHeader[:6] == "Basic " {
			// Extract the base64-encoded username and password
			username, password, ok := r.BasicAuth()

			// Verify the credentials here (e.g., check against a database)
			validCredentials := checkCredentials(username, password)

			// If the credentials are valid, proceed to the next handler
			if ok && validCredentials {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the credentials are invalid or missing, return a 401 Unauthorized status
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		w.WriteHeader(http.StatusUnauthorized)
	})
}

// checkCredentials is a dummy function to check if the provided credentials are valid
func checkCredentials(username string, password string) bool {
	// TODO: Replace with proper documentation
	// You can implement your own logic to validate the credentials
	// For demonstration purposes, we'll assume a valid username and password
	// validUsername := "admin"
	// validPassword := "password"

	return true
}
