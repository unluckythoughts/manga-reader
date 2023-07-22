package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetSourceMangas(ctx context.Context, id uint) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.WithContext(ctx).
		Joins("Source").
		Where("source_id = ?", id).
		Preload("Source").
		Limit(200).
		Find(&mangas).
		Error

	return mangas, err
}

// function to search the name by title
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

// function to search the name by title
func (r *Repository) SearchSourceMangasByTitle(ctx context.Context, source models.Source, query string) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.WithContext(ctx).
		Where("source_id = ? and LOWER(title) REGEXP ?", source.ID, strings.ToLower(query)).
		Preload("Source").
		Find(&mangas).
		Limit(100).
		Error

	return mangas, err
}

func (r *Repository) UpdateMangas(ctx context.Context, mangas *[]models.Manga) error {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "title"}, {Name: "source_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "synopsis", "image_url", "slug", "other_id"}),
		}, clause.Returning{}).
		Save(mangas).
		Error

	return err
}

func (r *Repository) GetManga(ctx context.Context, id uint) (models.Manga, error) {
	var manga models.Manga
	err := r.db.
		Where(id).
		Preload("Source").
		Preload("Chapters").
		Find(&manga).
		Error

	return manga, err
}

func (r *Repository) GetLatestMangaForSource(ctx context.Context, sourceID uint) (models.Manga, error) {
	var manga models.Manga
	err := r.db.
		Where("source_id  = ?", sourceID).
		Order("id DESC").
		Preload("Source").
		Preload("Chapters").
		First(&manga).
		Error

	return manga, err
}

func (r *Repository) GetSourceDomain(ctx context.Context, mangaID uint) (string, error) {
	var manga models.Manga
	err := r.db.
		Where(mangaID).
		Preload("Source").
		Find(&manga).
		Error

	if manga.Source.Domain != "" {
		return manga.Source.Domain, err
	}
	return "", errors.New("could not find domain")
}

func (r *Repository) GetChapter(ctx context.Context, id uint) (models.MangaChapter, error) {
	var chapter models.MangaChapter
	err := r.db.
		Where(id).
		Find(&chapter).
		Error

	return chapter, err
}

func (r *Repository) UpdateChapters(ctx context.Context, chapters *[]models.MangaChapter) error {
	err := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "number"}, {Name: "manga_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "title", "other_id"}),
		}, clause.Returning{}).
		Save(chapters).
		Error

	return err
}

func (r *Repository) UpdateChapterPages(ctx context.Context, id uint, pages []string) error {
	err := r.db.
		WithContext(ctx).
		Model(&models.MangaChapter{}).
		Where(id).
		UpdateColumn("image_urls", models.StrList(pages)).
		Error

	return err
}

func (r *Repository) DeleteChaptersBySource(ctx context.Context, sourceID uint) error {
	err := r.db.WithContext(ctx).
		Delete(&models.MangaChapter{}, "manga_id in (select id from manga where source_id=?)", sourceID).
		Error

	return err
}
