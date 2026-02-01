package service

import (
	"github.com/unluckythoughts/go-microservice/tools/auth"
)

// NOTE: User creation, updates, and authentication should be handled by auth.Service from go-microservice.
// These service methods are kept for admin/internal operations only.

// GetUsers retrieves a paginated list of users
func (s *ReaderService) GetUsers(page, limit int) ([]auth.User, int64, error) {
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
func (s *ReaderService) GetUserByID(id int) (*auth.User, error) {
	return s.db.GetUserByIDWithRelations(id, "Favorites")
}

// DeleteUser deletes a user by its ID (admin operation)
func (s *ReaderService) DeleteUser(id int) error {
	return s.db.DeleteUser(id)
}
