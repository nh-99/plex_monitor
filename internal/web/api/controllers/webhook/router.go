package webhook

import (
	"github.com/go-chi/chi"
)

// Routes returns the router for the webhook controller
func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", Entry)

	return router
}
