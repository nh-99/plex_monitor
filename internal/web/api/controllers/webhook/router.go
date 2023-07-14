package webhook

import (
	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", WebhookEntry)

	return router
}
