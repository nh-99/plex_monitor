package main

import (
	"compress/flate"
	"fmt"
	"net/http"
	"os"

	"plex_monitor/internal/config"
	"plex_monitor/internal/controllers/api/firehose"
	mediarequest "plex_monitor/internal/controllers/api/media_request"
	"plex_monitor/internal/controllers/api/user"
	"plex_monitor/internal/controllers/api/webhook"
	"plex_monitor/internal/controllers/interface/dashboard"
	user_interface "plex_monitor/internal/controllers/interface/user"
	"plex_monitor/internal/database"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

func routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
		logger.Logger("router", logrus.StandardLogger()), // Log API request calls
		middleware.Compress(flate.DefaultCompression),    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                       // Redirect slashes to no slash URL versions
		middleware.Recoverer,                             // Recover from panics without crashing server
	)

	router.Route("/api/v1", func(r chi.Router) {
		// Middleware for all the API routes
		r.Use(
			render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		)

		r.Mount("/firehose", firehose.Routes())
		r.Mount("/request", mediarequest.Routes())
		r.Mount("/users", user.Routes())
		r.Mount("/webhook", webhook.Routes())

		r.Get("/heartbeat", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}))
	})

	router.Route("/", func(r chi.Router) {
		r.Mount("/user", user_interface.Routes())
		r.Mount("/dashboard", dashboard.Routes())
	})

	return router
}

func initLogger() {
	conf := config.GetConfig()
	switch conf.Logging.Format {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(conf.Logging.ReportCaller)

	logrus.SetLevel(conf.Logging.GetLogrusLevel())
}

func main() {
	initLogger()
	router := routes()
	conf := config.GetConfig()
	database.InitDB(conf.Database.ConnectionString, conf.Database.Name)

	logrus.Info("Starting Plex Monitor Web...")

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logrus.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf.Server.Port), router))
}
