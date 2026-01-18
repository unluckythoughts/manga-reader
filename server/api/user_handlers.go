package api

import (
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/tools/web"
)

func (api *api) ListUsers(r web.Request) models.UsersPaginatedResponse {
	pageSize := r.GetURLParam("limit")
	pageNumber := r.GetURLParam("page")

	limit, err := strconv.Atoi(pageSize)
	if err != nil || limit <= 0 {
		limit = 20 // default limit
	}
	page, err := strconv.Atoi(pageNumber)
	if err != nil || page <= 0 {
		page = 1 // default page
	}

	users, total, err := api.s.GetUsers(page, limit)
	if err != nil {
		return models.UsersPaginatedResponse{
			Items: []models.User{},
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

	return models.UsersPaginatedResponse{
		Items:      users,
		Pagination: pagination,
	}
}

func (api *api) GetUser(r web.Request) (*models.User, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetUserByID(id)
}

func (api *api) CreateUser(r web.Request) (*models.User, error) {
	body := models.CreateUserRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.CreateUser(&body)
}

func (api *api) UpdateUser(r web.Request) (*models.User, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	body := models.UpdateUserRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.UpdateUser(id, &body)
}

func (api *api) DeleteUser(r web.Request) error {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return err
	}

	return api.s.DeleteUser(id)
}
