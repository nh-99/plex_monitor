package static

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func ViewCdn(w http.ResponseWriter, r *http.Request) {
	root := "./web/static"
	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(root)))
	fs.ServeHTTP(w, r)
}
