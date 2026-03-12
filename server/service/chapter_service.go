package service

import (
	"github.com/unluckythoughts/book-reader/server/connector"
	"github.com/unluckythoughts/book-reader/server/models"
)

// GetChapters retrieves a paginated list of chapters with optional book_id filter
func (s *ReaderService) GetChapters(page, limit int, bookID uint) ([]models.Chapter, int64, error) {
	offset := (page - 1) * limit

	chapters, err := s.db.ListChapters(offset, limit, bookID)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.db.CountChapters(bookID)
	if err != nil {
		return nil, 0, err
	}

	return chapters, total, nil
}

// GetChapterByID retrieves a chapter by its ID
func (s *ReaderService) GetChapterByID(id uint) (*models.Chapter, error) {
	return s.db.GetChapterByIDWithRelations(id, "Book")
}

// UpdateChapter updates an existing chapter
func (s *ReaderService) UpdateChapter(chapter *models.Chapter) error {
	sourceID := chapter.Book.SourceID
	src, err := s.db.GetSourceByID(sourceID)
	if err != nil {
		return err
	}

	conn, err := connector.GetConnector(src.Name)
	if err != nil {
		return err
	}

	// Use the connector to fetch the chapter details
	content, err := conn.GetChapterContent(chapter.URL)
	if err != nil {
		return err
	}
	chapter.Content = content

	// save chapter to database after fetching details
	return s.db.UpdateChapter(chapter)
}

// CreateChapter creates a new chapter
func (s *ReaderService) CreateChapter(input *models.CreateChapterRequest) (*models.Chapter, error) {
	chapter := &models.Chapter{
		URL:        input.URL,
		Title:      input.Title,
		BookID:     input.BookID,
		Number:     input.Number,
		Completed:  input.Completed,
		Downloaded: input.Downloaded,
		OtherID:    input.OtherID,
	}

	if err := s.db.CreateChapter(chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}
