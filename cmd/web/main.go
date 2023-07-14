package main

import (
	"compress/flate"
	"net/http"
	"os"

	"plex_monitor/internal/database"
	"plex_monitor/internal/web/api/controllers/firehose"
	"plex_monitor/internal/web/api/controllers/user"
	"plex_monitor/internal/web/api/controllers/webhook"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
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

		r.Mount("/firehose", firehose.Routes())
		r.Mount("/users", user.Routes())
		r.Mount("/webhook", webhook.Routes())

		r.Get("/heartbeat", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}))
	})

	return router
}

func initLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

func main() {
	router := Routes()
	database.InitDB(os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_NAME"))

	initLogger()

	logrus.Info("Starting Plex Monitor Web...")

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logrus.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	logrus.Fatal(http.ListenAndServe(":8080", router)) // Note, the port is usually gotten from the environment.
}
