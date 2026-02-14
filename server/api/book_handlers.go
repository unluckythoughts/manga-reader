package api

import (
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
)

func (api *api) ListBooks(r web.Request) (any, error) {
	pageSize := r.GetURLParam("limit")
	pageNumber := r.GetURLParam("page")
	sourceIDText := r.GetURLParam("source_id")

	limit, err := strconv.Atoi(pageSize)
	if err != nil || limit <= 0 {
		limit = 10 // default limit
	}
	page, err := strconv.Atoi(pageNumber)
	if err != nil || page <= 0 {
		page = 1 // default page
	}
	sourceID, err := strconv.Atoi(sourceIDText)
	if err != nil {
		sourceID = 0 // default sourceID
	}

	books, total, err := api.s.GetBooks(page, limit, uint(sourceID))
	if err != nil {
		return models.BooksPaginatedResponse{}, err
	}

	return models.BooksPaginatedResponse{
		Items: books,
		Pagination: models.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (api *api) GetBook(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetBookByID(uint(id))
}
