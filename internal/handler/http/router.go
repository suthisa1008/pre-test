package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/demo/product-api/docs"
)

func NewRouter(h *ProductHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Post("/product", h.CreateProduct)
	r.Patch("/product/{id}", h.PatchProduct)

	docs.SwaggerInfo.BasePath = "/"
	r.Get("/api-docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api-docs/index.html", http.StatusFound)
	})
	r.Get("/api-docs/*", httpSwagger.Handler(
		httpSwagger.URL("/api-docs/doc.json"),
	))

	return r
}
