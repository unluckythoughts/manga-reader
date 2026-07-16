package api

import (
	"github.com/unluckythoughts/book-reader/server/db"
	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/book-reader/server/service"
	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
	"github.com/unluckythoughts/go-microservice/v2/tools/worker"
	"github.com/unluckythoughts/go-microservice/v2/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type api struct {
	EnableAuth bool `env:"SERVICE_AUTH_ENABLE" envDefault:"false"`
	s          *service.ReaderService
	a          *auth.Service
}

// registerRoutes registers all API routes
func (a *api) registerRoutes(router web.Router) {

	// List sources route is public and does not require authentication
	// it retrieves a list of sources
	router.GET("/api/v1/public/sources", a.ListSources)

	// Get source route is public and does not require authentication
	// it retrieves a source by ID
	router.GET("/api/v1/public/sources/:id", a.GetSource)

	// List books route is public and does not require authentication in most cases,
	// it retrieves a paginated list of books
	// it supports filtering by book type, and source ID
	router.GET("/api/v1/public/books", a.ListBooks)

	// Get book route is public and does not require authentication
	// it retrieves a book by ID
	router.GET("/api/v1/public/books/:id", a.GetBook)

	// List chapters route is public and does not require authentication
	// it retrieves a paginated list of chapters for a book
	// it supports filtering by book ID
	router.GET("/api/v1/public/chapters", a.ListChapters)

	// Get chapter route is public and does not require authentication
	// it retrieves a chapter by ID
	router.GET("/api/v1/public/chapters/:id", a.GetChapter)

	if a.EnableAuth {
		auth.RegisterAuthRoutes(router, "/api/v1", a.a, models.UserRole)

		// Check if Google OAuth is configured
		if a.a.GoogleOauthConfig.ClientID != "" && a.a.GoogleOauthConfig.ClientSecret != "" {
			// OAUTH API
			router.POST("/api/v1/oauth/login/google", a.a.GoogleOAuthLogin)
		}

		// Protect Reader API routes
		router.UseFor("/api/v1/reader/", a.a.GetAuthMiddleware())
		
		// List sources route is protected and requires authentication
		// it retrieves a list of sources for the authenticated user
		router.GET("/api/v1/reader/sources", a.ListReaderSources)

		// List books route is protected and requires authentication
		// it retrieves a paginated list of books for the authenticated user
		// it supports filtering by search, book type, and source ID
		router.GET("/api/v1/reader/books", a.SearchReaderBooks)
		
		// Get Favorite route is protected and requires authentication
		// it retrieves all favorites for the authenticated user
		router.GET("/api/v1/reader/favorites", a.GetFavorite)

		// Create, Update, and Delete Favorite routes are protected and require authentication
		// they allow the authenticated user to create, update, and delete favorite sources and books
		// they support creating and updating favorites with progress tracking
		router.POST("/api/v1/reader/favorites", a.CreateFavorite)
		router.PUT("/api/v1/reader/favorites", a.UpdateFavorite)
		router.PATCH("/api/v1/reader/favorites", a.UpdateFavoriteProgress)
		router.DELETE("/api/v1/reader/favorites", a.DeleteFavorite)
	}
}

func Register(router web.Router, gormDB *gorm.DB, l *zap.Logger, w *worker.Worker) {
	s := service.New(db.New(gormDB), w, l)
	api := &api{s: s}

	// Initialize auth service (needed even when auth is disabled for dummy user)
	a := auth.New(auth.Options{
		DB:                gormDB,
		Logger:            l.Named("auth"),
		TokenValidInHours: 24,
	})
	api.a = a

	// Check if auth should be enabled
	utils.ParseEnvironmentVars(api)

	api.registerRoutes(router)
}
