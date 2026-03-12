package service

import (
	"github.com/unluckythoughts/book-reader/server/connector"
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

// UpdateBook retrieves a book by its url from the connector itself
func (s *ReaderService) UpdateBook(book *models.Book) error {
	if book.Source.Domain == "" {
		source, err := s.db.GetSourceByID(book.SourceID)
		if err != nil {
			return err
		}
		book.Source = source
	}

	conn, err := connector.GetConnector(book.Source.Name)
	if err != nil {
		return err
	}

	// Use the connector to fetch the book details
	chapters, err := conn.GetBookChapters(book.URL)
	if err != nil {
		return err
	}
	book.Chapters = chapters

	synopsis, err := conn.GetBookSynopsis(book.URL)
	if err != nil {
		return err
	}
	book.Synopsis = synopsis

	// save book to database after fetching details
	err = s.db.UpdateBook(book)
	if err != nil {
		return err
	}

	return nil
}
