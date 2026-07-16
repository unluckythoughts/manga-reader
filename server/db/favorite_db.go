package db

import (
	"errors"
	"time"

	"github.com/unluckythoughts/book-reader/server/models"
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

func (d *DB) UpdateFavorite(favorite *models.Favorite) error {
	favorite.UpdatedAt = time.Now()
	if err := d.db.Save(favorite).Error; err != nil {
		return err
	}
	return nil
}

// GetFavoriteByID retrieves a favorite by its ID
func (d *DB) GetFavoriteByID(id uint) (*models.Favorite, error) {
	var favorite models.Favorite

	query := d.db.Where("id = ?", id)

	if err := query.First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("favorite not found")
		}
		return nil, err
	}
	return &favorite, nil
}

// GetFavoriteByUserID retrieves a favorite by its user ID
func (d *DB) GetFavoriteByUserID(userID uint) (*models.Favorite, error) {
	var favorite models.Favorite

	query := d.db.Where("user_id = ?", userID)

	if err := query.First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("favorite not found")
		}
		return nil, err
	}
	return &favorite, nil
}

// DeleteFavoriteByUserID soft deletes a favorite by user_id
func (d *DB) DeleteFavoriteByUserID(userID uint) error {
	result := d.db.Where("user_id = ?", userID).Delete(&models.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// DeleteFavorite deletes a favorite by its ID
func (d *DB) DeleteFavorite(id uint) error {
	if err := d.db.Where("id = ?", id).Delete(&models.Favorite{}).Error; err != nil {
		return err
	}
	return nil
}

// HardDeleteFavoriteByUserID permanently deletes a favorite from the database
func (d *DB) HardDeleteFavoriteByUserID(userID uint) error {
	result := d.db.Unscoped().Where("user_id = ?", userID).Delete(&models.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// HardDeleteFavorite permanently deletes a favorite from the database
func (d *DB) HardDeleteFavorite(id uint) error {
	result := d.db.Unscoped().Delete(&models.Favorite{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// FavoriteExists checks if a favorite exists for a user and book
func (d *DB) FavoriteExists(userID uint) (bool, error) {
	var count int64
	if err := d.db.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
