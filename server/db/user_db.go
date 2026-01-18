package db

import (
	"errors"
	"time"

	"github.com/unluckythoughts/book-reader/server/models"
	"gorm.io/gorm"
)

// CreateUser creates a new user in the database
func (d *DB) CreateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	if err := d.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by its ID
func (d *DB) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := d.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByIDWithRelations retrieves a user by its ID with related data
func (d *DB) GetUserByIDWithRelations(id int, preload ...string) (*models.User, error) {
	var user models.User
	query := d.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByName retrieves a user by its name
func (d *DB) GetUserByName(name string) (*models.User, error) {
	var user models.User
	if err := d.db.Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (d *DB) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	if err := d.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUserFields updates specific fields of a user
func (d *DB) UpdateUserFields(id int, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()
	result := d.db.Model(&models.User{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// DeleteUser soft deletes a user by setting DeletedAt
func (d *DB) DeleteUser(id int) error {
	result := d.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// HardDeleteUser permanently deletes a user from the database
func (d *DB) HardDeleteUser(id int) error {
	result := d.db.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// ListUsers retrieves all users with pagination
func (d *DB) ListUsers(offset, limit int) ([]models.User, error) {
	var users []models.User
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("name ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// ListUsersWithRelations retrieves all users with related data
func (d *DB) ListUsersWithRelations(offset, limit int, preload ...string) ([]models.User, error) {
	var users []models.User
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.Order("name ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CountUsers returns the total number of users
func (d *DB) CountUsers() (int64, error) {
	var count int64
	if err := d.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAllUsers retrieves all users without pagination
func (d *DB) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := d.db.Order("name ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// SearchUsers searches users by name
func (d *DB) SearchUsers(searchTerm string, offset, limit int) ([]models.User, error) {
	var users []models.User
	query := d.db.Where("name ILIKE ?", "%"+searchTerm+"%").Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("name ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
