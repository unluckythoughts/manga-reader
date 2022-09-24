package reader

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/reader/api"
	"github.com/unluckythoughts/manga-reader/reader/service"
)

func RegisterRoutes(router web.Router, service *service.Service) {
	h := api.New(service)

	router.GET("/api/:type/source", h.SourceListHandler)
	router.POST("/api/:type/source", h.SourceItemListHandler)
	router.POST("/api/:type/source/item", h.SourceItemHandler)
	router.POST("/api/:type/source/chapter", h.SourceItemChapterHandler)

	router.POST("/api/:type/source/search", h.SourceItemSearchHandler)

	router.GET("/api/:type/library", h.GetFavoriteListHandler)
	router.POST("/api/:type/library", h.AddFavoriteHandler)
	router.PATCH("/api/:type/library", h.UpdateAllFavoriteHandler)
	router.GET("/api/:type/library/:favoriteID/update", h.UpdateFavoriteHandler)
	router.DELETE("/api/:type/library/:favoriteID/remove", h.DelFavoriteHandler)
	router.PUT("/api/:type/library/:favoriteID/chapter/:chapterID/progress/:pageID", h.UpdateFavoriteProgressHandler)
}
