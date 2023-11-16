package firehose

import (
	"plex_monitor/internal/config"
	"plex_monitor/internal/controllers/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

// Routes returns the router for the firehose endpoints
func Routes() *chi.Mux {
	router := chi.NewRouter()
	globals := config.GetGlobals()
	tokenAuth := globals.JWTAuth

	// Protected endpoints
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Custom middleware for X to add user to request context, for easy access
		r.Use(middleware.CreateUserContext)

		// Private endpoints
		r.Get("/", Firehose)
	})

	return router
}
