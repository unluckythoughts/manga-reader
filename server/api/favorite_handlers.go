package api

import (
	"net/http"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
)

func (api *api) ListFavorites(r web.Request) (any, error) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		return nil, web.NewError(http.StatusUnauthorized, err)
	}

	fav, err := api.s.GetFavoriteByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return fav, nil
}

func (api *api) GetFavorite(r web.Request) (any, error) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		return nil, web.NewError(http.StatusUnauthorized, err)
	}

	return api.s.GetFavoriteByUserID(user.ID)
}

func (api *api) CreateFavorite(r web.Request) (any, error) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		return nil, web.NewError(http.StatusUnauthorized, err)
	}

	body := models.AddFavoriteRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.CreateFavorite(&body, user.ID)
}

func (api *api) UpdateFavorite(r web.Request) (any, error) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		return nil, web.NewError(http.StatusUnauthorized, err)
	}

	body := models.AddFavoriteRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.UpdateFavorite(user.ID, &body)
}

func (api *api) UpdateFavoriteProgress(r web.Request) (any, error) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		return nil, web.NewError(http.StatusUnauthorized, err)
	}

	body := models.UpdateFavoriteProgressRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	var p models.Progress
	p.ScanValue(body.Chapter, body.Level)

	updates := models.AddFavoriteRequest{
		BookID:   body.BookID,
		Progress: string(p),
	}

	return api.s.UpdateFavorite(user.ID, &updates)
}

func (api *api) DeleteFavorite(r web.Request) (any, error) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		return nil, web.NewError(http.StatusUnauthorized, err)
	}

	return nil, api.s.DeleteFavoriteByUserID(user.ID)
}
