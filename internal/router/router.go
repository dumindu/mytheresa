package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"github.com/mytheresa/go-hiring-challenge/internal/app/product"
	"github.com/mytheresa/go-hiring-challenge/internal/router/middleware"
	reqL "github.com/mytheresa/go-hiring-challenge/internal/router/middleware/requestlog"
	"github.com/mytheresa/go-hiring-challenge/internal/util/logger"
)

func New(db *gorm.DB, l *logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("."))
	})

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "pragma"},
		AllowCredentials: true,
		ExposedHeaders:   []string{},
		MaxAge:           300,
	})

	r.Route("/v1", func(r chi.Router) {
		r.Use(cors.Handler)
		r.Use(middleware.RequestID)
		r.Use(middleware.ContentTypeJSON)

		productAPI := product.New(db, l)
		r.Method(http.MethodGet, "/products", reqL.NewHandler(productAPI.GetAll, l))
	})

	return r
}
