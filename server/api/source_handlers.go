package api

import (
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/tools/web"
)

func (api *api) ListSources(r web.Request) models.SourcesPaginatedResponse {
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
		return models.SourcesPaginatedResponse{
			Items: []models.Source{},
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

	return models.SourcesPaginatedResponse{
		Items:      sources,
		Pagination: pagination,
	}
}

func (api *api) GetSource(r web.Request) (*models.Source, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetSourceByID(id)
}
