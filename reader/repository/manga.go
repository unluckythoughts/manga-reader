package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetSourceMangas(ctx context.Context, domain string) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.WithContext(ctx).
		Joins("Source").
		Where("Source.domain = ?", domain).
		Preload("Source").
		Limit(200).
		Find(&mangas).
		Error

	return mangas, err
}

func (r *Repository) SearchMangasByTitle(ctx context.Context, query string) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.WithContext(ctx).
		Where("LOWER(title) REGEXP ?", strings.ToLower(query)).
		Preload("Source").
		Find(&mangas).
		Limit(100).
		Error

	return mangas, err
}

func (r *Repository) UpdateMangas(ctx context.Context, mangas *[]models.Manga) error {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			DoUpdates: clause.AssignmentColumns([]string{"synopsis", "image_url", "slug", "other_id"}),
		}, clause.Returning{}).
		Save(mangas).
		Error

	return err
}

func (r *Repository) UpdateMangaByName(ctx context.Context, manga *models.Manga) error {
	dbManga := models.Manga{}
	err := r.db.WithContext(ctx).
		First(&dbManga, "source_id=? and title=?", manga.SourceID, manga.Title).
		Error
	if err != nil {
		return err
	} else if dbManga.ID <= 0 {
		return errors.New("could find manga to update")
	}

	manga.ID = dbManga.ID
	return r.db.WithContext(ctx).Save(manga).Error
}

func (r *Repository) GetManga(ctx context.Context, url string) (models.Manga, error) {
	var manga models.Manga
	err := r.db.
		Where(&models.Manga{URL: url}).
		Preload("Source").
		Preload("Chapters").
		Find(&manga).
		Error

	return manga, err
}

func (r *Repository) UpdateChapters(ctx context.Context, chapters *[]models.MangaChapter) error {
	err := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			DoUpdates: clause.AssignmentColumns([]string{"number", "title", "other_id"}),
		}, clause.Returning{}).
		Save(chapters).
		Error

	return err
}

func (r *Repository) DeleteChaptersBySource(ctx context.Context, sourceID int) error {
	err := r.db.WithContext(ctx).
		Delete(&models.MangaChapter{}, "manga_id in (select id from manga where source_id=?)", sourceID).
		Error

	return err
}
