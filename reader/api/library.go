package api

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func (h *Handlers) AddFavoriteHandler(r web.Request) (interface{}, error) {
	body := models.SourceRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, web.BadRequest(err)
	}

	if isMangaRequest(r) {
		return h.s.AddMangaFavorite(r.GetContext(), body.URL)
	} else if isNovelRequest(r) {
		return h.s.AddNovelFavorite(r.GetContext(), body.URL)
	}

	return nil, web.BadRequest()
}

func (h *Handlers) GetFavoriteListHandler(r web.Request) (interface{}, error) {
	if isMangaRequest(r) {
		return h.s.GetMangaFavorites(r.GetContext())
	} else if isNovelRequest(r) {
		return h.s.GetNovelFavorites(r.GetContext())
	}

	return nil, web.BadRequest()
}

func (h *Handlers) DelFavoriteHandler(r web.Request) (interface{}, error) {
	favoriteID := r.GetRouteParam("favoriteID")

	id, err := strconv.Atoi(favoriteID)
	if err != nil {
		return nil, web.BadRequest(err)
	}

	if isMangaRequest(r) {
		return nil, h.s.DelMangaFavorite(r.GetContext(), id)
	} else if isNovelRequest(r) {
		return nil, h.s.DelNovelFavorite(r.GetContext(), id)
	}

	return nil, web.BadRequest()
}

func (h *Handlers) UpdateFavoriteHandler(r web.Request) (interface{}, error) {
	favoriteID := r.GetRouteParam("favoriteID")

	id, err := strconv.Atoi(favoriteID)
	if err != nil {
		return nil, web.BadRequest(err)
	}

	if isMangaRequest(r) {
		return h.s.UpdateMangaFavorite(r.GetContext(), id)
	} else if isNovelRequest(r) {
		return h.s.UpdateNovelFavorite(r.GetContext(), id)
	}

	return nil, web.BadRequest()

}

func (h *Handlers) UpdateAllFavoriteHandler(r web.Request) (interface{}, error) {
	if isMangaRequest(r) {
		return nil, h.s.UpdateAllMangaFavorite(r.GetContext())
	} else if isNovelRequest(r) {
		return nil, h.s.UpdateAllNovelFavorite(r.GetContext())
	}

	return nil, web.BadRequest()
}

func (h *Handlers) UpdateFavoriteProgressHandler(r web.Request) (interface{}, error) {
	strFavoriteID := r.GetRouteParam("favoriteID")
	strChapterID := r.GetRouteParam("chapterID")
	strPageID := r.GetRouteParam("pageID")

	favoriteID, err := strconv.Atoi(strFavoriteID)
	if err != nil {
		return nil, web.BadRequest(errors.Wrapf(err, "could get favoriteID route param"))
	}

	chapterID, err := strconv.ParseFloat(strChapterID, 64)
	if err != nil {
		return nil, web.BadRequest(errors.Wrapf(err, "could get chapterID route param"))
	}

	pageID, err := strconv.ParseFloat(strPageID, 64)
	if err != nil {
		pageID = 0
	}

	progress := models.StrFloatList{chapterID, pageID}
	if isMangaRequest(r) {
		return nil, h.s.UpdateMangaFavoriteProgress(r.GetContext(), favoriteID, progress)
	} else if isNovelRequest(r) {
		return nil, h.s.UpdateNovelFavoriteProgress(r.GetContext(), favoriteID, progress)
	}

	return nil, web.BadRequest()
}
