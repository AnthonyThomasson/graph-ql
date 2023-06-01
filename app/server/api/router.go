package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

func getRouter(db *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/query", getGqlHandler(db))
	r.Get("/docs", getPlaygroundHandler())
	setStaticHandlers(r)
	return r
}
