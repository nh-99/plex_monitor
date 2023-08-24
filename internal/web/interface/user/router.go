package user

import (
	"os"
	"plex_monitor/internal/web/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

// Routes creates a REST router for the user API
func Routes() *chi.Mux {
	router := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)

	// Protected endpoints
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		// Custom middleware for plex_monitor to add user to request context, for easy access
		r.Use(middleware.CreateUserContext)

		// TODO: add routes here
	})

	// Public endpoints
	router.Group(func(r chi.Router) {
		// Custom middleware for plex_monitor to add user to request context, for easy access
		r.Use(middleware.CreateUserContext)

		r.Get("/login", ViewLogin)
		r.Post("/login", PerformLogin)
	})

	return router
}
