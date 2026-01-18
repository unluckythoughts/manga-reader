package client

import (
	"fmt"

	"github.com/unluckythoughts/manga-reader/server/models"
)

// ListSources retrieves a paginated list of sources
func (c *Client) ListSources(page, limit int) (*models.SourcesPaginatedResponse, error) {
	url := fmt.Sprintf("/api/v1/sources?page=%d&limit=%d", page, limit)

	var response models.SourcesPaginatedResponse
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetSource retrieves a source by ID
func (c *Client) GetSource(id int) (*models.Source, error) {
	url := fmt.Sprintf("/api/v1/sources/%d", id)

	var response models.Source
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
