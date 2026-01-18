package db

import (
	"errors"
	"time"

	"github.com/unluckythoughts/manga-reader/server/models"
	"gorm.io/gorm"
)

// CreateBook creates a new book in the database
func (d *DB) CreateBook(book *models.Book) error {
	book.UpdatedAt = time.Now()
	if err := d.db.Create(book).Error; err != nil {
		return err
	}
	return nil
}

// GetBookByID retrieves a book by its ID
func (d *DB) GetBookByID(id int) (*models.Book, error) {
	var book models.Book
	if err := d.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

// GetBookByIDWithRelations retrieves a book by its ID with related data
func (d *DB) GetBookByIDWithRelations(id int, preload ...string) (*models.Book, error) {
	var book models.Book
	query := d.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

// GetBookByURL retrieves a book by its URL
func (d *DB) GetBookByURL(url string) (*models.Book, error) {
	var book models.Book
	if err := d.db.Where("url = ?", url).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

// UpdateBook updates an existing book
func (d *DB) UpdateBook(book *models.Book) error {
	book.UpdatedAt = time.Now()
	if err := d.db.Save(book).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBook soft deletes a book by setting DeletedAt
func (d *DB) DeleteBook(id int) error {
	result := d.db.Delete(&models.Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}

// HardDeleteBook permanently deletes a book from the database
func (d *DB) HardDeleteBook(id int) error {
	result := d.db.Unscoped().Delete(&models.Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}

// ListBooks retrieves all books with optional filters and returns total count
func (d *DB) ListBooks(offset, limit int, bookType string, sourceID int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	// Build base query for counting
	countQuery := d.db.Model(&models.Book{})
	if bookType != "" {
		countQuery = countQuery.Where("type = ?", bookType)
	}
	if sourceID > 0 {
		countQuery = countQuery.Where("source_id = ?", sourceID)
	}

	// Get total count
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build query for fetching records
	query := d.db.Offset(offset)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if bookType != "" {
		query = query.Where("type = ?", bookType)
	}
	if sourceID > 0 {
		query = query.Where("source_id = ?", sourceID)
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, 0, err
	}
	return books, total, nil
}

// ListBooksWithRelations retrieves all books with related data
func (d *DB) ListBooksWithRelations(offset, limit int, bookType string, sourceID int, preload ...string) ([]models.Book, error) {
	var books []models.Book
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if bookType != "" {
		query = query.Where("type = ?", bookType)
	}

	if sourceID > 0 {
		query = query.Where("source_id = ?", sourceID)
	}

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// CountBooks returns the total number of books
func (d *DB) CountBooks(bookType string, sourceID int) (int64, error) {
	var count int64
	query := d.db.Model(&models.Book{})

	if bookType != "" {
		query = query.Where("type = ?", bookType)
	}

	if sourceID > 0 {
		query = query.Where("source_id = ?", sourceID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SearchBooks searches books by title
func (d *DB) SearchBooks(searchTerm string, offset, limit int) ([]models.Book, error) {
	var books []models.Book
	query := d.db.Where("title ILIKE ?", "%"+searchTerm+"%").Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// GetBookBySourceAndOtherID retrieves a book by source_id and other_id
func (d *DB) GetBookBySourceAndOtherID(sourceID int, otherID string) (*models.Book, error) {
	var book models.Book
	if err := d.db.Where("source_id = ? AND other_id = ?", sourceID, otherID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}
