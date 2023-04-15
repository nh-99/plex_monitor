package dashboard

import (
	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()

	// Public endpoints
	router.Group(func(r chi.Router) {
		r.Get("/", ViewDashboard)
	})

	return router
}
