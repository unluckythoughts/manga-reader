package api

import (
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
)

func (api *api) ListChapters(r web.Request) (any, error) {
	pageSize := r.GetURLParam("limit")
	pageNumber := r.GetURLParam("page")
	bookIDStr := r.GetURLParam("book_id")

	limit, err := strconv.Atoi(pageSize)
	if err != nil || limit <= 0 {
		limit = 20 // default limit
	}
	page, err := strconv.Atoi(pageNumber)
	if err != nil || page <= 0 {
		page = 1 // default page
	}

	var bookID int
	if bookIDStr != "" {
		if id, err := strconv.Atoi(bookIDStr); err == nil {
			bookID = id
		}
	}

	chapters, total, err := api.s.GetChapters(page, limit, bookID)
	if err != nil {
		return models.ChaptersPaginatedResponse{}, err
	}

	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	pagination.CalculateTotalPages()

	return models.ChaptersPaginatedResponse{
		Items:      chapters,
		Pagination: pagination,
	}, nil
}

func (api *api) GetChapter(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetChapterByID(id)
}
