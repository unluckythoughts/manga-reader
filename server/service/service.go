package service

import (
	"github.com/unluckythoughts/manga-reader/server/db"
	"github.com/unluckythoughts/manga-reader/server/models"
)

type ReaderService struct {
	db *db.DB
}

// New creates a new instance of ReaderService
func New(db *db.DB) *ReaderService {
	return &ReaderService{db: db}
}

// GetBooks retrieves a paginated list of books with optional filters
func (s *ReaderService) GetBooks(page, limit, sourceID int) ([]models.Book, int64, error) {
	offset := (page - 1) * limit

	return s.db.ListBooks(offset, limit, "", sourceID)
}

// GetBookByID retrieves a book by its ID
func (s *ReaderService) GetBookByID(id int) (*models.Book, error) {
	return s.db.GetBookByIDWithRelations(id, "Source", "Chapters")
}
