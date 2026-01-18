package client

import (
	"fmt"

	"github.com/unluckythoughts/manga-reader/server/models"
)

// ListBooks retrieves a paginated list of books
func (c *Client) ListBooks(page, limit, sourceID int) (*models.BooksPaginatedResponse, error) {
	url := fmt.Sprintf("/api/v1/books?page=%d&limit=%d", page, limit)
	if sourceID > 0 {
		url += fmt.Sprintf("&source_id=%d", sourceID)
	}

	var response models.BooksPaginatedResponse
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBook retrieves a book by ID
func (c *Client) GetBook(id int) (*models.Book, error) {
	url := fmt.Sprintf("/api/v1/books/%d", id)

	var response models.Book
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
