package api

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/server/db"
	"github.com/unluckythoughts/manga-reader/server/service"
	"gorm.io/gorm"
)

type api struct {
	s *service.ReaderService
}

// registerRoutes registers all API routes
func (a *api) registerRoutes(router web.Router) {
	// Books API
	router.GET("/api/v1/books", a.ListBooks)
	router.GET("/api/v1/books/:id", a.GetBook)

	// Sources API
	router.GET("/api/v1/sources", a.ListSources)
	router.GET("/api/v1/sources/:id", a.GetSource)
}

func Register(router web.Router, gormDB *gorm.DB) {
	s := service.New(db.New(gormDB))
	api := &api{s: s}

	api.registerRoutes(router)
}
