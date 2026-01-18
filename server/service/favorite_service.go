package service

import (
	"time"

	"github.com/unluckythoughts/manga-reader/server/models"
)

// GetFavorites retrieves a paginated list of favorites with optional filters
func (s *ReaderService) GetFavorites(page, limit, userID int) ([]models.Favorite, int64, error) {
	offset := (page - 1) * limit

	favorites, err := s.db.ListFavorites(offset, limit, userID)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.db.CountFavorites(userID)
	if err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// GetFavoriteByID retrieves a favorite by its ID
func (s *ReaderService) GetFavoriteByID(id int) (*models.Favorite, error) {
	return s.db.GetFavoriteByIDWithRelations(id, "User", "Book")
}

// CreateFavorite creates a new favorite
func (s *ReaderService) CreateFavorite(input *models.CreateFavoriteRequest) (*models.Favorite, error) {
	favorite := &models.Favorite{
		UserID:    input.UserID,
		BookID:    input.BookID,
		UpdatedAt: time.Now(),
	}

	if err := s.db.CreateFavorite(favorite); err != nil {
		return nil, err
	}

	return favorite, nil
}

// UpdateFavorite updates an existing favorite
func (s *ReaderService) UpdateFavorite(id int, input *models.UpdateFavoriteRequest) (*models.Favorite, error) {
	favorite, err := s.db.GetFavoriteByID(id)
	if err != nil {
		return nil, err
	}

	favorite.UpdatedAt = time.Now()

	if err := s.db.UpdateFavorite(favorite); err != nil {
		return nil, err
	}

	return favorite, nil
}

// DeleteFavorite deletes a favorite by its ID
func (s *ReaderService) DeleteFavorite(id int) error {
	return s.db.DeleteFavorite(id)
}
