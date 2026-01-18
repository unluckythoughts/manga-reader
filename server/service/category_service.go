package service

import (
	"time"

	"github.com/unluckythoughts/manga-reader/server/models"
)

// GetCategories retrieves a paginated list of categories
func (s *ReaderService) GetCategories(page, limit int) ([]models.Category, int64, error) {
	offset := (page - 1) * limit

	categories, err := s.db.ListCategories(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.db.CountCategories()
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// GetCategoryByID retrieves a category by its ID
func (s *ReaderService) GetCategoryByID(id int) (*models.Category, error) {
	return s.db.GetCategoryByID(id)
}

// CreateCategory creates a new category
func (s *ReaderService) CreateCategory(input *models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:      input.Name,
		UpdatedAt: time.Now(),
	}

	if err := s.db.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates an existing category
func (s *ReaderService) UpdateCategory(id int, input *models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.db.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = input.Name
	category.UpdatedAt = time.Now()

	if err := s.db.UpdateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category by its ID
func (s *ReaderService) DeleteCategory(id int) error {
	return s.db.DeleteCategory(id)
}
