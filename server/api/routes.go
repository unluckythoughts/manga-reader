package api

import (
	"github.com/unluckythoughts/book-reader/server/db"
	"github.com/unluckythoughts/book-reader/server/service"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
	"github.com/unluckythoughts/go-microservice/v2/tools/worker"
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

	// Chapters API
	router.GET("/api/v1/chapters", a.ListChapters)
	router.GET("/api/v1/chapters/:id", a.GetChapter)

	// Sources API
	router.GET("/api/v1/sources", a.ListSources)
	router.GET("/api/v1/sources/:id", a.GetSource)

	// Users API (admin operations only - for auth, use auth package handlers)
	// Authentication routes (login, register, etc.) should be handled by auth.Service
	router.GET("/api/v1/users", a.ListUsers)
	router.GET("/api/v1/users/:id", a.GetUser)

	// Categories API
	router.GET("/api/v1/categories", a.ListCategories)
	router.GET("/api/v1/categories/:id", a.GetCategory)
	router.POST("/api/v1/categories", a.CreateCategory)
	router.PUT("/api/v1/categories/:id", a.UpdateCategory)
	router.DELETE("/api/v1/categories/:id", a.DeleteCategory)

	// Favorites API
	router.GET("/api/v1/favorites", a.ListFavorites)
	router.GET("/api/v1/favorites/:id", a.GetFavorite)
	router.POST("/api/v1/favorites", a.CreateFavorite)
	router.PUT("/api/v1/favorites/:id", a.UpdateFavorite)
	router.PATCH("/api/v1/favorites/:id", a.UpdateFavoriteProgress)
	router.DELETE("/api/v1/favorites/:id", a.DeleteFavorite)
}

func Register(router web.Router, gormDB *gorm.DB, w *worker.Worker) {
	s := service.New(db.New(gormDB), w)
	api := &api{s: s}

	api.registerRoutes(router)
}
