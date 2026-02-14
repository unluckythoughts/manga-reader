package client

import (
	"fmt"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
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
func (c *Client) GetUser(id uint) (*auth.User, error) {
	url := fmt.Sprintf("/api/v1/users/%d", id)

	var response auth.User
	_, err := c.client.GetResponse(url, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateUser creates a new user
// Note: This endpoint has been removed. Use auth.Service.GetRoleUserRegister instead.
// This method is kept for backwards compatibility but will not work without auth routes.
func (c *Client) CreateUser(request *auth.RegisterRequest) (*auth.User, error) {
	url := "/api/v1/users"

	var response auth.User
	_, err := c.client.PostResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateUser updates an existing user
// Note: This endpoint has been removed. Use auth.Service.UpdateUserHandler instead.
// This method is kept for backwards compatibility but will not work without auth routes.
func (c *Client) UpdateUser(id uint, request *auth.UpdateUserRequest) (*auth.User, error) {
	url := fmt.Sprintf("/api/v1/users/%d", id)

	var response auth.User
	_, err := c.client.PutResponse(url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteUser deletes a user by ID
func (c *Client) DeleteUser(id uint) error {
	url := fmt.Sprintf("/api/v1/users/%d", id)

	_, err := c.client.DeleteResponse(url, nil, nil)
	return err
}
