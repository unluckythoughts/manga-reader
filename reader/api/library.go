package api

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func (h *Handlers) AddFavoriteHandler(r web.Request) (interface{}, error) {
	body := models.SourceMangaRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, err
	}

	return h.s.AddFavorite(r.GetContext(), body.MangaURL)
}

func (h *Handlers) GetFavoriteListHandler(r web.Request) (interface{}, error) {
	return h.s.GetFavorites(r.GetContext())
}

func (h *Handlers) DelFavoriteHandler(r web.Request) (interface{}, error) {
	favoriteID := r.GetRouteParam("favoriteID")

	id, err := strconv.Atoi(favoriteID)
	if err != nil {
		return nil, err
	}

	return nil, h.s.DelFavorite(r.GetContext(), id)
}

func (h *Handlers) UpdateFavoriteHandler(r web.Request) (interface{}, error) {
	favoriteID := r.GetRouteParam("favoriteID")

	id, err := strconv.Atoi(favoriteID)
	if err != nil {
		return nil, err
	}

	return h.s.UpdateFavorite(r.GetContext(), id)
}

func (h *Handlers) UpdateAllFavoriteHandler(r web.Request) (interface{}, error) {
	err := h.s.UpdateAllFavorite(r.GetContext())
	return nil, err
}

func (h *Handlers) UpdateFavoriteProgressHandler(r web.Request) (interface{}, error) {
	strFavoriteID := r.GetRouteParam("favoriteID")
	strChapterID := r.GetRouteParam("chapterID")
	strPageID := r.GetRouteParam("pageID")

	favoriteID, err := strconv.Atoi(strFavoriteID)
	if err != nil {
		return nil, errors.Wrapf(err, "could get favoriteID route param")
	}

	chapterID, err := strconv.Atoi(strChapterID)
	if err != nil {
		return nil, errors.Wrapf(err, "could get chapterID route param")
	}

	pageID, err := strconv.Atoi(strPageID)
	if err != nil {
		pageID = 0
	}

	progress := models.StrIntList{chapterID, pageID}
	return nil, h.s.UpdateFavoriteProgress(r.GetContext(), favoriteID, progress)
}
