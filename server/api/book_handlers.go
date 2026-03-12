package api

import (
	"net/http"
	"strconv"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
)

func (api *api) ListBooks(r web.Request) (any, error) {
	pageSize := r.GetURLParam("limit")
	pageNumber := r.GetURLParam("page")
	sourceIDText := r.GetURLParam("source_id")
	bookTypeText := r.GetURLParam("book_type")
	search := r.GetURLParam("search")

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
	bookType := "" // default book type (no filter)
	if bookTypeText == string(models.BookTypeManga) ||
		bookTypeText == string(models.BookTypeNovel) {
		bookType = bookTypeText
	}

	books, total, err := api.s.GetBooks(page, limit, search, bookType, uint(sourceID))
	if err != nil {
		return models.BooksPaginatedResponse{}, err
	}

	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	pagination.CalculateTotalPages()

	return models.BooksPaginatedResponse{
		Items:      books,
		Pagination: pagination,
	}, nil
}

func (api *api) GetBook(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, web.NewError(http.StatusNotFound, err)
	}

	forceText := r.GetURLParam("force")
	force := false
	if forceText != "" {
		force, err = strconv.ParseBool(forceText)
		if err != nil {
			return nil, web.NewError(http.StatusUnprocessableEntity, err)
		}
	}

	book, err := api.s.GetBookByID(uint(id))
	if err != nil {
		return nil, err
	}

	if len(book.Chapters) > 0 || book.Synopsis != "" {
		if !force {
			// if force is not enabled return db data
			return book, nil
		}
	}

	// If the book has no chapters and an empty synopsis,
	// fetch the details from the source and update
	err = api.s.UpdateBook(book)
	if err != nil {
		return nil, err
	}

	return book, nil
}
