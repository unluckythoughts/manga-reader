package service

import (
	"time"

	"github.com/unluckythoughts/manga-reader/server/models"
)

// GetUsers retrieves a paginated list of users
func (s *ReaderService) GetUsers(page, limit int) ([]models.User, int64, error) {
	offset := (page - 1) * limit

	users, err := s.db.ListUsers(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.db.CountUsers()
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID retrieves a user by its ID
func (s *ReaderService) GetUserByID(id int) (*models.User, error) {
	return s.db.GetUserByIDWithRelations(id, "Favorites")
}

// CreateUser creates a new user
func (s *ReaderService) CreateUser(input *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Name:      input.Name,
		UpdatedAt: time.Now(),
	}

	if err := s.db.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *ReaderService) UpdateUser(id int, input *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.db.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = input.Name
	user.UpdatedAt = time.Now()

	if err := s.db.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by its ID
func (s *ReaderService) DeleteUser(id int) error {
	return s.db.DeleteUser(id)
}
