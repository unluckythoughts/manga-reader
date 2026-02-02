package db

import (
	"errors"
	"time"

	"github.com/unluckythoughts/book-reader/server/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateChapter creates a new chapter in the database
func (d *DB) CreateChapter(chapter *models.Chapter) error {
	chapter.UpdatedAt = time.Now()
	if err := d.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "title", "book_id", "number", "content", "upload_date", "completed", "downloaded", "other_id", "updated_at"}),
		},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "book_id"}, {Name: "number"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "title", "content", "upload_date", "completed", "downloaded", "other_id", "updated_at"}),
		},
	).Create(chapter).Error; err != nil {
		return err
	}
	return nil
}

// CreateChaptersBatch creates multiple chapters in a single transaction
func (d *DB) CreateChaptersBatch(chapters []models.Chapter) error {
	now := time.Now()
	for i := range chapters {
		chapters[i].UpdatedAt = now
	}

	if err := d.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "title", "book_id", "number", "content", "upload_date", "completed", "downloaded", "other_id", "updated_at"}),
		},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "book_id"}, {Name: "number"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "title", "content", "upload_date", "completed", "downloaded", "other_id", "updated_at"}),
		},
	).Create(&chapters).Error; err != nil {
		return err
	}
	return nil
}

// GetChapterByID retrieves a chapter by its ID
func (d *DB) GetChapterByID(id int) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := d.db.First(&chapter, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		return nil, err
	}
	return &chapter, nil
}

// GetChapterByIDWithRelations retrieves a chapter by its ID with related data
func (d *DB) GetChapterByIDWithRelations(id int, preload ...string) (*models.Chapter, error) {
	var chapter models.Chapter
	query := d.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(&chapter, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		return nil, err
	}
	return &chapter, nil
}

// GetChapterByURL retrieves a chapter by its URL
func (d *DB) GetChapterByURL(url string) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := d.db.Where("url = ?", url).First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		return nil, err
	}
	return &chapter, nil
}

// GetChaptersByBookID retrieves all chapters for a specific book
func (d *DB) GetChaptersByBookID(bookID int, offset, limit int) ([]models.Chapter, error) {
	var chapters []models.Chapter
	query := d.db.Where("book_id = ?", bookID).Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("number ASC").Find(&chapters).Error; err != nil {
		return nil, err
	}
	return chapters, nil
}

// UpdateChapter updates an existing chapter
func (d *DB) UpdateChapter(chapter *models.Chapter) error {
	chapter.UpdatedAt = time.Now()
	if err := d.db.Save(chapter).Error; err != nil {
		return err
	}
	return nil
}

// UpdateChapterFields updates specific fields of a chapter
func (d *DB) UpdateChapterFields(id int, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()
	result := d.db.Model(&models.Chapter{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("chapter not found")
	}
	return nil
}

// MarkChapterAsCompleted marks a chapter as completed
func (d *DB) MarkChapterAsCompleted(id int) error {
	return d.UpdateChapterFields(id, map[string]interface{}{"completed": true})
}

// MarkChapterAsDownloaded marks a chapter as downloaded
func (d *DB) MarkChapterAsDownloaded(id int) error {
	return d.UpdateChapterFields(id, map[string]interface{}{"downloaded": true})
}

// DeleteChapter soft deletes a chapter by setting DeletedAt
func (d *DB) DeleteChapter(id int) error {
	result := d.db.Delete(&models.Chapter{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("chapter not found")
	}
	return nil
}

// DeleteChaptersByBookID deletes all chapters for a specific book
func (d *DB) DeleteChaptersByBookID(bookID int) error {
	if err := d.db.Where("book_id = ?", bookID).Delete(&models.Chapter{}).Error; err != nil {
		return err
	}
	return nil
}

// HardDeleteChapter permanently deletes a chapter from the database
func (d *DB) HardDeleteChapter(id int) error {
	result := d.db.Unscoped().Delete(&models.Chapter{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("chapter not found")
	}
	return nil
}

// ListChapters retrieves all chapters with optional filters
func (d *DB) ListChapters(offset, limit int, bookID int) ([]models.Chapter, error) {
	var chapters []models.Chapter
	query := d.db.Offset(offset)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if bookID > 0 {
		query = query.Where("book_id = ?", bookID)
	}

	if err := query.Order("number ASC").Find(&chapters).Error; err != nil {
		return nil, err
	}
	return chapters, nil
}

// CountChapters returns the total number of chapters
func (d *DB) CountChapters(bookID int) (int64, error) {
	var count int64
	query := d.db.Model(&models.Chapter{})

	if bookID > 0 {
		query = query.Where("book_id = ?", bookID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetChapterByBookAndNumber retrieves a chapter by book_id and chapter number
func (d *DB) GetChapterByBookAndNumber(bookID int, number string) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := d.db.Where("book_id = ? AND number = ?", bookID, number).First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		return nil, err
	}
	return &chapter, nil
}

// GetLatestChapterByBookID retrieves the latest chapter for a book
func (d *DB) GetLatestChapterByBookID(bookID int) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := d.db.Where("book_id = ?", bookID).Order("upload_date DESC").First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		return nil, err
	}
	return &chapter, nil
}
