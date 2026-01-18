package client

import (
	"fmt"

	"github.com/unluckythoughts/book-reader/server/models"
)

// ListUsers retrieves a paginated list of users
func (c *Client) ListUsers(page, limit int) (*models.UsersPaginatedResponse, error) {
	url := fmt.Sprintf("/api/v1/users?page=%d&limit=%d", page, limit)

	var response models.UsersPaginatedResponse
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUser retrieves a user by ID
func (c *Client) GetUser(id int) (*models.User, error) {
	url := fmt.Sprintf("/api/v1/users/%d", id)

	var response models.User
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateUser creates a new user
func (c *Client) CreateUser(request *models.CreateUserRequest) (*models.User, error) {
	url := "/api/v1/users"

	var response models.User
	_, err := c.client.PostResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateUser updates an existing user
func (c *Client) UpdateUser(id int, request *models.UpdateUserRequest) (*models.User, error) {
	url := fmt.Sprintf("/api/v1/users/%d", id)

	var response models.User
	_, err := c.client.PutResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteUser deletes a user by ID
func (c *Client) DeleteUser(id int) error {
	url := fmt.Sprintf("/api/v1/users/%d", id)

	_, err := c.client.DeleteResponse(url, nil, nil)
	return err
}
