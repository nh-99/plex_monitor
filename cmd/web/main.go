package main

import (
	"compress/flate"
	"log"
	"net/http"
	"os"

	"plex_monitor/internal/database"
	"plex_monitor/internal/web/api/controllers/user"
	"plex_monitor/internal/web/interface/controllers/dashboard"
	"plex_monitor/internal/web/interface/controllers/static"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger, // Log API request calls
		middleware.Compress(flate.DefaultCompression), // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
	)

	router.Route("/api/v1", func(r chi.Router) {
		// Middleware for all the API routes
		r.Use(
			render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		)

		r.Mount("/users", user.Routes())
	})

	router.Route("/", func(r chi.Router) {
		r.Mount("/", dashboard.Routes())
		r.Mount("/static", static.Routes())
	})

	return router
}

func main() {
	router := Routes()
	models.InitDB(os.Getenv("DATABASE_URL"))

	log.Printf("Starting Plex Monitor Web...")

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Fatal(http.ListenAndServe(":8080", router)) // Note, the port is usually gotten from the environment.
}
