package db

import (
	"errors"
	"time"

	"github.com/unluckythoughts/manga-reader/server/models"
	"gorm.io/gorm"
)

// CreateFavorite creates a new favorite in the database
func (d *DB) CreateFavorite(favorite *models.Favorite) error {
	favorite.UpdatedAt = time.Now()
	if err := d.db.Create(favorite).Error; err != nil {
		return err
	}
	return nil
}

// GetFavoriteByID retrieves a favorite by its ID
func (d *DB) GetFavoriteByID(id int) (*models.Favorite, error) {
	var favorite models.Favorite
	if err := d.db.First(&favorite, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("favorite not found")
		}
		return nil, err
	}
	return &favorite, nil
}

// GetFavoriteByIDWithRelations retrieves a favorite by its ID with related data
func (d *DB) GetFavoriteByIDWithRelations(id int, preload ...string) (*models.Favorite, error) {
	var favorite models.Favorite
	query := d.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(&favorite, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("favorite not found")
		}
		return nil, err
	}
	return &favorite, nil
}

// GetFavoriteByUserAndBook retrieves a favorite by user_id and book_id
func (d *DB) GetFavoriteByUserAndBook(userID, bookID int) (*models.Favorite, error) {
	var favorite models.Favorite
	if err := d.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("favorite not found")
		}
		return nil, err
	}
	return &favorite, nil
}

// GetFavoritesByUserID retrieves all favorites for a specific user
func (d *DB) GetFavoritesByUserID(userID int, offset, limit int) ([]models.Favorite, error) {
	var favorites []models.Favorite
	query := d.db.Where("user_id = ?", userID).Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("updated_at DESC").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

// GetFavoritesByUserIDWithRelations retrieves all favorites for a user with related data
func (d *DB) GetFavoritesByUserIDWithRelations(userID int, offset, limit int, preload ...string) ([]models.Favorite, error) {
	var favorites []models.Favorite
	query := d.db.Where("user_id = ?", userID).Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.Order("updated_at DESC").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

// GetFavoritesByBookID retrieves all favorites for a specific book
func (d *DB) GetFavoritesByBookID(bookID int) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := d.db.Where("book_id = ?", bookID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

// UpdateFavorite updates an existing favorite
func (d *DB) UpdateFavorite(favorite *models.Favorite) error {
	favorite.UpdatedAt = time.Now()
	if err := d.db.Save(favorite).Error; err != nil {
		return err
	}
	return nil
}

// UpdateFavoriteFields updates specific fields of a favorite
func (d *DB) UpdateFavoriteFields(id int, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()
	result := d.db.Model(&models.Favorite{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// UpdateFavoriteProgress updates the progress field of a favorite
func (d *DB) UpdateFavoriteProgress(id int, progress string) error {
	return d.UpdateFavoriteFields(id, map[string]interface{}{"progress": progress})
}

// UpdateFavoriteCategories updates the categories field of a favorite
func (d *DB) UpdateFavoriteCategories(id int, categories string) error {
	return d.UpdateFavoriteFields(id, map[string]interface{}{"categories": categories})
}

// DeleteFavorite soft deletes a favorite by setting DeletedAt
func (d *DB) DeleteFavorite(id int) error {
	result := d.db.Delete(&models.Favorite{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// DeleteFavoriteByUserAndBook deletes a favorite by user_id and book_id
func (d *DB) DeleteFavoriteByUserAndBook(userID, bookID int) error {
	result := d.db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// DeleteFavoritesByUserID deletes all favorites for a specific user
func (d *DB) DeleteFavoritesByUserID(userID int) error {
	if err := d.db.Where("user_id = ?", userID).Delete(&models.Favorite{}).Error; err != nil {
		return err
	}
	return nil
}

// HardDeleteFavorite permanently deletes a favorite from the database
func (d *DB) HardDeleteFavorite(id int) error {
	result := d.db.Unscoped().Delete(&models.Favorite{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// ListFavorites retrieves all favorites with optional filters
func (d *DB) ListFavorites(offset, limit int, userID, bookID int) ([]models.Favorite, error) {
	var favorites []models.Favorite
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if bookID > 0 {
		query = query.Where("book_id = ?", bookID)
	}

	if err := query.Order("updated_at DESC").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

// CountFavorites returns the total number of favorites
func (d *DB) CountFavorites(userID, bookID int) (int64, error) {
	var count int64
	query := d.db.Model(&models.Favorite{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if bookID > 0 {
		query = query.Where("book_id = ?", bookID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FavoriteExists checks if a favorite exists for a user and book
func (d *DB) FavoriteExists(userID, bookID int) (bool, error) {
	var count int64
	if err := d.db.Model(&models.Favorite{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
