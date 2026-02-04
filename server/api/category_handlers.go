package api

import (
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
)

func (api *api) ListCategories(r web.Request) (any, error) {
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

	categories, total, err := api.s.GetCategories(page, limit)
	if err != nil {
		return models.CategoriesPaginatedResponse{}, err
	}

	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	pagination.CalculateTotalPages()

	return models.CategoriesPaginatedResponse{
		Items:      categories,
		Pagination: pagination,
	}, nil
}

func (api *api) GetCategory(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetCategoryByID(id)
}

func (api *api) CreateCategory(r web.Request) (any, error) {
	body := models.CreateCategoryRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.CreateCategory(&body)
}

func (api *api) UpdateCategory(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	body := models.UpdateCategoryRequest{}
	if err := r.GetValidatedBody(&body); err != nil {
		return nil, err
	}

	return api.s.UpdateCategory(id, &body)
}

func (api *api) DeleteCategory(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return nil, api.s.DeleteCategory(id)
}
