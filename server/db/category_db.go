package db

import (
	"errors"
	"time"

	"github.com/unluckythoughts/book-reader/server/models"
	"gorm.io/gorm"
)

// CreateCategory creates a new category in the database
func (d *DB) CreateCategory(category *models.Category) error {
	category.UpdatedAt = time.Now()
	if err := d.db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

// GetCategoryByID retrieves a category by its ID
func (d *DB) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := d.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// GetCategoryByName retrieves a category by its name
func (d *DB) GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	if err := d.db.Where("name = ?", name).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category
func (d *DB) UpdateCategory(category *models.Category) error {
	category.UpdatedAt = time.Now()
	if err := d.db.Save(category).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCategoryFields updates specific fields of a category
func (d *DB) UpdateCategoryFields(id uint, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()
	result := d.db.Model(&models.Category{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// DeleteCategory soft deletes a category by setting DeletedAt
func (d *DB) DeleteCategory(id uint) error {
	result := d.db.Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// HardDeleteCategory permanently deletes a category from the database
func (d *DB) HardDeleteCategory(id uint) error {
	result := d.db.Unscoped().Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// ListCategories retrieves all categories with pagination
func (d *DB) ListCategories(offset, limit int) ([]models.Category, error) {
	var categories []models.Category
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// CountCategories returns the total number of categories
func (d *DB) CountCategories() (int64, error) {
	var count int64
	if err := d.db.Model(&models.Category{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAllCategories retrieves all categories without pagination
func (d *DB) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := d.db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// SearchCategories searches categories by name
func (d *DB) SearchCategories(searchTerm string, offset, limit int) ([]models.Category, error) {
	var categories []models.Category
	query := d.db.Where("name ILIKE ?", "%"+searchTerm+"%").Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// CategoryExists checks if a category exists by name
func (d *DB) CategoryExists(name string) (bool, error) {
	var count int64
	if err := d.db.Model(&models.Category{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
