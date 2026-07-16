package service

import (
	"maps"
	"slices"

	"github.com/unluckythoughts/book-reader/server/models"
)

// GetFavoriteByUserID retrieves a favorite by its user ID
func (s *ReaderService) GetFavoriteByUserID(user_id uint) (*models.FavoriteResponse, error) {
	fav, err := s.db.GetFavoriteByUserID(user_id)
	if err != nil {
		return nil, err
	}

	sources, err := s.db.GetSourcesByIDs(fav.Data.Sources)
	if err != nil {
		return nil, err
	}

	books, err := s.db.GetBooksByIDs(slices.Collect(maps.Keys(fav.Data.Books)))
	if err != nil {
		return nil, err
	}

	bookMap := make(map[uint]models.FavoriteBookValue)

	for _, book := range books {
		progress, exists := fav.Data.Books[book.ID]
		if !exists {
			progress = "0:0" // default progress if not found
		}
		bookMap[book.ID] = models.FavoriteBookValue{
			Book:     book,
			Progress: progress,
		}
	}

	favResp := &models.FavoriteResponse{
		Sources:  sources,
		Books:    bookMap,
		Settings: fav.Data.Settings,
	}

	return favResp, nil
}

// CreateFavorite creates a new favorite
func (s *ReaderService) CreateFavorite(input *models.AddFavoriteRequest, userID uint) (*models.Favorite, error) {
	fav, err := s.db.GetFavoriteByUserID(userID)
	if err == nil {
		return fav, nil
	}

	if fav == nil {
		fav = &models.Favorite{
			UserID: userID,
			Data: models.FavoriteData{
				Sources:  []uint{},
				Books:    make(map[uint]models.Progress),
				Settings: make(map[string]string),
			},
		}
	}

	if input.SourceID != 0 {
		fav.Data.Sources = append(fav.Data.Sources, input.SourceID)
	}

	if input.BookID != 0 {
		if input.Progress == "" {
			input.Progress = "0:0" // default progress
		}

		fav.Data.Books[input.BookID] = models.Progress(input.Progress)
	}

	if err := s.db.CreateFavorite(fav); err != nil {
		return nil, err
	}

	return fav, nil
}

// UpdateFavorite updates an existing favorite
func (s *ReaderService) UpdateFavorite(userID uint, input *models.AddFavoriteRequest) (*models.Favorite, error) {
	favorite, err := s.db.GetFavoriteByUserID(userID)
	if err != nil {
		return nil, err
	}
	if input.SourceID != 0 {
		favorite.Data.Sources = append(favorite.Data.Sources, input.SourceID)
	}

	if input.BookID != 0 {
		if input.Progress == "" {
			input.Progress = "0:0" // default progress
		}
		favorite.Data.Books[input.BookID] = models.Progress(input.Progress)
	}

	if err := s.db.UpdateFavorite(favorite); err != nil {
		return nil, err
	}

	return favorite, nil
}

// DeleteFavorite deletes a favorite by its ID
func (s *ReaderService) DeleteFavorite(id uint) error {
	return s.db.DeleteFavorite(id)
}

// DeleteFavoriteByUserID deletes a favorite by its user ID
func (s *ReaderService) DeleteFavoriteByUserID(userID uint) error {
	return s.db.DeleteFavoriteByUserID(userID)
}
