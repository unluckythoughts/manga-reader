package reader

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/reader/api"
	"github.com/unluckythoughts/manga-reader/reader/service"
)

func RegisterRoutes(router web.Router, service *service.Service) {
	h := api.New(service)

	router.GET("/proxy/:url", h.ProxyHandler)

	router.GET("/source", h.SourceListHandler)
	router.POST("/source", h.SourceMangaListHandler)
	router.POST("/source/manga", h.SourceMangaHandler)
	router.POST("/source/chapter", h.SourceMangaChapterHandler)

	router.POST("/source/search", h.SourceMangaSearchHandler)

	router.GET("/library", h.GetFavoriteListHandler)
	router.POST("/library", h.AddFavoriteHandler)
	router.PATCH("/library", h.UpdateAllFavoriteHandler)
	router.GET("/library/:favoriteID/update", h.UpdateFavoriteHandler)
	router.DELETE("/library/:favoriteID/remove", h.DelFavoriteHandler)
	router.PUT("/library/:favoriteID/chapter/:chapterID/progress/:pageID", h.UpdateFavoriteProgressHandler)
}
