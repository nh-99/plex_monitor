package user

import (
	"plex_monitor/internal/web/middleware"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

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

		// Custom middleware for X to add user to request context, for easy access
		r.Use(middleware.CreateUserContext)
	})

	// Public endpoints
	router.Group(func(r chi.Router) {
		// Custom middleware for X to add user to request context, for easy access
		r.Use(middleware.CreateUserContext)

		r.Post("/login", PerformLogin)
	})

	return router
}
