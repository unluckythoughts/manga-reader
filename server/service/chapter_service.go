package service

import (
	"time"

	"github.com/unluckythoughts/manga-reader/server/models"
)

// GetChapters retrieves a paginated list of chapters with optional book_id filter
func (s *ReaderService) GetChapters(page, limit, bookID int) ([]models.Chapter, int64, error) {
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
func (s *ReaderService) GetChapterByID(id int) (*models.Chapter, error) {
	return s.db.GetChapterByIDWithRelations(id, "Book")
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
		UpdatedAt:  time.Now(),
	}

	if err := s.db.CreateChapter(chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}
