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

	// Users API
	router.GET("/api/v1/users", a.ListUsers)
	router.GET("/api/v1/users/:id", a.GetUser)
	router.POST("/api/v1/users", a.CreateUser)
	router.PUT("/api/v1/users/:id", a.UpdateUser)
	router.DELETE("/api/v1/users/:id", a.DeleteUser)

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
	router.DELETE("/api/v1/favorites/:id", a.DeleteFavorite)
}

func Register(router web.Router, gormDB *gorm.DB) {
	s := service.New(db.New(gormDB))
	api := &api{s: s}

	api.registerRoutes(router)
}
