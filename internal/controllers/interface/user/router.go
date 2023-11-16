package user

import (
	"plex_monitor/internal/config"
	"plex_monitor/internal/controllers/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

// Routes creates a REST router for the user API
func Routes() *chi.Mux {
	router := chi.NewRouter()
	globals := config.GetGlobals()
	tokenAuth := globals.JWTAuth

	// Protected endpoints
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Custom middleware for plex_monitor to add user to request context, for easy access
		r.Use(middleware.CreateUserContext)

		// TODO: add routes here
	})

	// Public endpoints
	router.Group(func(r chi.Router) {
		r.Get("/login", ViewLogin)
		r.Post("/login", PerformLogin)
	})

	return router
}
