package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
)

func (api *api) ListSources(r web.Request) (any, error) {
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

	sources, total, err := api.s.GetSources(page, limit)
	if err != nil {
		return models.SourcesPaginatedResponse{}, err
	}

	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	pagination.CalculateTotalPages()

	return models.SourcesPaginatedResponse{
		Items:      sources,
		Pagination: pagination,
	}, nil
}

func (api *api) GetSource(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetSourceByID(uint(id))
}

func (api *api) ListReaderSources(r web.Request) (any, error) {
	pageSize := r.GetURLParam("limit")
	pageNumber := r.GetURLParam("page")

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil || user == nil {
		return nil, web.NewError(http.StatusUnauthorized, fmt.Errorf("unauthorized: %+v", err))
	}

	limit, err := strconv.Atoi(pageSize)
	if err != nil || limit <= 0 {
		limit = 20 // default limit
	}
	page, err := strconv.Atoi(pageNumber)
	if err != nil || page <= 0 {
		page = 1 // default page
	}

	sources, total, err := api.s.GetSources(page, limit)
	if err != nil {
		return models.SourcesPaginatedResponse{}, err
	}

	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	pagination.CalculateTotalPages()

	return models.SourcesPaginatedResponse{
		Items:      sources,
		Pagination: pagination,
	}, nil
}

