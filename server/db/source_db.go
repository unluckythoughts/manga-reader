package db

import (
	"errors"
	"time"

	"github.com/unluckythoughts/book-reader/server/models"
	"gorm.io/gorm"
)

// CreateSource creates a new source in the database
func (d *DB) CreateSource(source *models.Source) error {
	source.UpdatedAt = time.Now()
	if err := d.db.Create(source).Error; err != nil {
		return err
	}
	return nil
}

// GetSourceByID retrieves a source by its ID
func (d *DB) GetSourceByID(id int) (*models.Source, error) {
	var source models.Source
	if err := d.db.First(&source, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("source not found")
		}
		return nil, err
	}
	return &source, nil
}

// GetSourceByIDWithRelations retrieves a source by its ID with related data
func (d *DB) GetSourceByIDWithRelations(id int, preload ...string) (*models.Source, error) {
	var source models.Source
	query := d.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(&source, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("source not found")
		}
		return nil, err
	}
	return &source, nil
}

// GetSourceByName retrieves a source by its name
func (d *DB) GetSourceByName(name string) (*models.Source, error) {
	var source models.Source
	if err := d.db.Where("name = ?", name).First(&source).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("source not found")
		}
		return nil, err
	}
	return &source, nil
}

// GetSourceByDomain retrieves a source by its domain
func (d *DB) GetSourceByDomain(domain string) (*models.Source, error) {
	var source models.Source
	if err := d.db.Where("domain = ?", domain).First(&source).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("source not found")
		}
		return nil, err
	}
	return &source, nil
}

// UpdateSource updates an existing source
func (d *DB) UpdateSource(source *models.Source) error {
	source.UpdatedAt = time.Now()
	if err := d.db.Save(source).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSourceFields updates specific fields of a source
func (d *DB) UpdateSourceFields(id int, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()
	result := d.db.Model(&models.Source{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("source not found")
	}
	return nil
}

// DeleteSource soft deletes a source by setting DeletedAt
func (d *DB) DeleteSource(id int) error {
	result := d.db.Delete(&models.Source{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("source not found")
	}
	return nil
}

// HardDeleteSource permanently deletes a source from the database
func (d *DB) HardDeleteSource(id int) error {
	result := d.db.Unscoped().Delete(&models.Source{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("source not found")
	}
	return nil
}

// ListSources retrieves all sources with pagination
func (d *DB) ListSources(offset, limit int) ([]models.Source, error) {
	var sources []models.Source
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("name ASC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

// ListSourcesWithRelations retrieves all sources with related data
func (d *DB) ListSourcesWithRelations(offset, limit int, preload ...string) ([]models.Source, error) {
	var sources []models.Source
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.Order("name ASC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

// CountSources returns the total number of sources
func (d *DB) CountSources() (int64, error) {
	var count int64
	if err := d.db.Model(&models.Source{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAllSources retrieves all sources without pagination
func (d *DB) GetAllSources() ([]models.Source, error) {
	var sources []models.Source
	if err := d.db.Order("name ASC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

// SearchSources searches sources by name or domain
func (d *DB) SearchSources(searchTerm string, offset, limit int) ([]models.Source, error) {
	var sources []models.Source
	query := d.db.Where("name ILIKE ? OR domain ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%").Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("name ASC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}
