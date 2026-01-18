package service

import "github.com/unluckythoughts/manga-reader/server/models"

// GetSources retrieves a paginated list of sources
func (s *ReaderService) GetSources(page, limit int) ([]models.Source, int64, error) {
	offset := (page - 1) * limit

	sources, err := s.db.ListSources(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.db.CountSources()
	if err != nil {
		return nil, 0, err
	}

	return sources, total, nil
}

// GetSourceByID retrieves a source by its ID
func (s *ReaderService) GetSourceByID(id int) (*models.Source, error) {
	return s.db.GetSourceByID(id)
}
