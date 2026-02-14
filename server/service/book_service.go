package service

import (
	"github.com/unluckythoughts/book-reader/server/models"
)

// GetBooks retrieves a paginated list of books with optional filters
func (s *ReaderService) GetBooks(page, limit int, sourceID uint) ([]models.Book, int64, error) {
	offset := (page - 1) * limit

	return s.db.ListBooks(offset, limit, "", sourceID)
}

// GetBookByID retrieves a book by its ID
func (s *ReaderService) GetBookByID(id uint) (*models.Book, error) {
	return s.db.GetBookByIDWithRelations(id, "Source", "Chapters")
}
