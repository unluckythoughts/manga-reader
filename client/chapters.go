package client

import (
	"fmt"

	"github.com/unluckythoughts/manga-reader/server/models"
)

// ListChapters retrieves a paginated list of chapters
func (c *Client) ListChapters(page, limit, bookID int) (*models.ChaptersPaginatedResponse, error) {
	url := fmt.Sprintf("/api/v1/chapters?page=%d&limit=%d", page, limit)
	if bookID > 0 {
		url += fmt.Sprintf("&book_id=%d", bookID)
	}

	var response models.ChaptersPaginatedResponse
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetChapter retrieves a chapter by ID
func (c *Client) GetChapter(id int) (*models.Chapter, error) {
	url := fmt.Sprintf("/api/v1/chapters/%d", id)

	var response models.Chapter
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
