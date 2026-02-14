package client

import (
	"fmt"

	"github.com/unluckythoughts/book-reader/server/models"
)

// ListCategories retrieves a paginated list of categories
func (c *Client) ListCategories(page, limit int) (*models.CategoriesPaginatedResponse, error) {
	url := fmt.Sprintf("/api/v1/categories?page=%d&limit=%d", page, limit)

	var response models.CategoriesPaginatedResponse
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCategory retrieves a category by ID
func (c *Client) GetCategory(id uint) (*models.Category, error) {
	url := fmt.Sprintf("/api/v1/categories/%d", id)

	var response models.Category
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateCategory creates a new category
func (c *Client) CreateCategory(request *models.CreateCategoryRequest) (*models.Category, error) {
	url := "/api/v1/categories"

	var response models.Category
	_, err := c.client.PostResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateCategory updates an existing category
func (c *Client) UpdateCategory(id uint, request *models.UpdateCategoryRequest) (*models.Category, error) {
	url := fmt.Sprintf("/api/v1/categories/%d", id)

	var response models.Category
	_, err := c.client.PutResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteCategory deletes a category by ID
func (c *Client) DeleteCategory(id uint) error {
	url := fmt.Sprintf("/api/v1/categories/%d", id)

	_, err := c.client.DeleteResponse(url, nil, nil)
	return err
}
