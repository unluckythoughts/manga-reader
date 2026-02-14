package client

import (
	"fmt"

	"github.com/unluckythoughts/book-reader/server/models"
)

// ListFavorites retrieves a paginated list of favorites
func (c *Client) ListFavorites(page, limit int, userID, bookID uint) (*models.FavoritesPaginatedResponse, error) {
	url := fmt.Sprintf("/api/v1/favorites?page=%d&limit=%d", page, limit)
	if userID > 0 {
		url += fmt.Sprintf("&user_id=%d", userID)
	}
	if bookID > 0 {
		url += fmt.Sprintf("&book_id=%d", bookID)
	}

	var response models.FavoritesPaginatedResponse
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetFavorite retrieves a favorite by ID
func (c *Client) GetFavorite(id uint) (*models.Favorite, error) {
	url := fmt.Sprintf("/api/v1/favorites/%d", id)

	var response models.Favorite
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateFavorite creates a new favorite
func (c *Client) CreateFavorite(request *models.CreateFavoriteRequest) (*models.Favorite, error) {
	url := "/api/v1/favorites"

	var response models.Favorite
	_, err := c.client.PostResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateFavorite updates an existing favorite
func (c *Client) UpdateFavorite(id uint, request *models.UpdateFavoriteRequest) (*models.Favorite, error) {
	url := fmt.Sprintf("/api/v1/favorites/%d", id)

	var response models.Favorite
	_, err := c.client.PutResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateFavoriteProgress updates the reading progress of a favorite
func (c *Client) UpdateFavoriteProgress(id uint, request *models.UpdateFavoriteProgressRequest) (*models.Favorite, error) {
	url := fmt.Sprintf("/api/v1/favorites/%d", id)

	var response models.Favorite
	_, err := c.client.PatchResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteFavorite deletes a favorite by ID
func (c *Client) DeleteFavorite(id uint) error {
	url := fmt.Sprintf("/api/v1/favorites/%d", id)

	_, err := c.client.DeleteResponse(url, nil, nil)
	return err
}
