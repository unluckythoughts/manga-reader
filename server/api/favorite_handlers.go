package api

import (
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/tools/web"
)

func (api *api) ListFavorites(r web.Request) models.FavoritesPaginatedResponse {
	pageSize := r.GetURLParam("limit")
	pageNumber := r.GetURLParam("page")
	userIDStr := r.GetURLParam("user_id")

	limit, err := strconv.Atoi(pageSize)
	if err != nil || limit <= 0 {
		limit = 20 // default limit
	}
	page, err := strconv.Atoi(pageNumber)
	if err != nil || page <= 0 {
		page = 1 // default page
	}

	var userID int
	if userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = id
		}
	}

	favorites, total, err := api.s.GetFavorites(page, limit, userID)
	if err != nil {
		return models.FavoritesPaginatedResponse{
			Items: []models.Favorite{},
			Pagination: models.Pagination{
				Page:  page,
				Limit: limit,
				Total: 0,
			},
		}
	}

	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	pagination.CalculateTotalPages()

	return models.FavoritesPaginatedResponse{
		Items:      favorites,
		Pagination: pagination,
	}
}

func (api *api) GetFavorite(r web.Request) (*models.Favorite, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetFavoriteByID(id)
}

func (api *api) CreateFavorite(r web.Request) (*models.Favorite, error) {
	body := models.CreateFavoriteRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.CreateFavorite(&body)
}

func (api *api) UpdateFavorite(r web.Request) (*models.Favorite, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	body := models.UpdateFavoriteRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.UpdateFavorite(id, &body)
}

func (api *api) UpdateFavoriteProgress(r web.Request) (*models.Favorite, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	body := models.UpdateFavoriteProgressRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	progress := models.NewList([]string{
		strconv.Itoa(body.Chapter),
		strconv.Itoa(body.Level),
	})
	updatedBody := models.UpdateFavoriteRequest{
		Progress: progress.String(),
	}

	return api.s.UpdateFavorite(id, &updatedBody)
}

func (api *api) DeleteFavorite(r web.Request) error {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return err
	}

	return api.s.DeleteFavorite(id)
}
