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
	a          *auth.Auth
}

// registerRoutes registers all API routes
func (a *api) registerRoutes(router web.Router) {
	if a.EnableAuth {
		// Auth API
		router.POST("/api/v1/auth/login", a.a.LoginHandler)
		router.POST("/api/v1/auth/register", a.a.GetRegisterHandlerForUserRole(models.UserRole))
		router.POST("/api/v1/auth/logout", a.a.LogoutHandler)

		// Check if Google OAuth is configured
		if a.a.GoogleOauthConfig.ClientID != "" && a.a.GoogleOauthConfig.ClientSecret != "" {
			// OAUTH API
			router.POST("/api/v1/oauth/login/google", a.a.GoogleOAuthLogin)
		}

		// Verification API
		router.PATCH("/api/v1/auth/verify/:target", a.a.SendTokenHandler)
		router.GET("/api/v1/auth/verify/:target/:token", a.a.VerifyTokenHandler)

		// User API
		router.GET("/api/v1/user", a.a.GetUserHandler)
		router.PUT("/api/v1/user", a.a.UpdateUserHandler)
		router.PATCH("/api/v1/user/change-password", a.a.ChangePasswordHandler)
		router.GET("/api/v1/user/reset-password", a.a.ResetPasswordHandler)
		router.PATCH("/api/v1/user/update-password", a.a.UpdatePasswordHandler)

		// Protect Reader API routes
		router.UseFor("/api/v1/reader/", a.a.GetAuthMiddleware())
	} else {
		// If auth is disabled, create a dummy user and set it in the context for all requests
		a.a.CreateUser(&auth.User{
			Name:     "dummy",
			Email:    "dummy@example.com",
			Password: "dummy",          // In a real application, use a secure password and hash it
			Role:     models.AdminRole, // Assign admin role for testing purposes
		})
	}

	// Books API
	router.GET("/api/v1/reader/books", a.ListBooks)
	router.GET("/api/v1/reader/books/:id", a.GetBook)

	// Chapters API
	router.GET("/api/v1/reader/chapters", a.ListChapters)
	router.GET("/api/v1/reader/chapters/:id", a.GetChapter)

	// Sources API
	router.GET("/api/v1/reader/sources", a.ListSources)
	router.GET("/api/v1/reader/sources/:id", a.GetSource)

	// Categories API
	router.GET("/api/v1/reader/categories", a.ListCategories)
	router.GET("/api/v1/reader/categories/:id", a.GetCategory)
	router.POST("/api/v1/reader/categories", a.CreateCategory)
	router.PUT("/api/v1/reader/categories/:id", a.UpdateCategory)
	router.DELETE("/api/v1/reader/categories/:id", a.DeleteCategory)

	// Favorites API
	router.GET("/api/v1/reader/favorites", a.ListFavorites)
	router.GET("/api/v1/reader/favorites/:id", a.GetFavorite)
	router.POST("/api/v1/reader/favorites", a.CreateFavorite)
	router.PUT("/api/v1/reader/favorites/:id", a.UpdateFavorite)
	router.PATCH("/api/v1/reader/favorites/:id", a.UpdateFavoriteProgress)
	router.DELETE("/api/v1/reader/favorites/:id", a.DeleteFavorite)
}

func Register(router web.Router, gormDB *gorm.DB, l *zap.Logger, w *worker.Worker) {
	s := service.New(db.New(gormDB), w)
	api := &api{s: s}

	// enable auth if configured
	utils.ParseEnvironmentVars(api)
	if api.EnableAuth {
		a := auth.NewAuthService(auth.Options{
			DB:                gormDB,
			Logger:            l.Named("auth"),
			TokenValidInHours: 24,
		})
		api.a = a
	}

	api.registerRoutes(router)
}
