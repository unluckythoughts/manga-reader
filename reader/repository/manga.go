package repository

import (
	"context"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetSourceMangas(ctx context.Context, domain string) ([]models.Manga, error) {
	var mangas []models.Manga
	err := r.db.WithContext(ctx).
		Joins("Source").
		Where("Source.domain = ?", domain).
		Find(&mangas).
		Error

	return mangas, err
}

func (r *Repository) UpdateMangas(ctx context.Context, mangas *[]models.Manga) error {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			UpdateAll: false,
			DoNothing: true,
		}, clause.Returning{}).
		Save(mangas).
		Error

	return err
}

func (r *Repository) GetManga(ctx context.Context, url string) (models.Manga, error) {
	var manga models.Manga
	err := r.db.
		Where(&models.Manga{URL: url}).
		Preload("Chapters").
		Find(&manga).
		Error

	return manga, err
}

func (r *Repository) UpdateChapters(ctx context.Context, chapters *[]models.Chapter) error {
	err := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			UpdateAll: false,
			DoNothing: true,
		}, clause.Returning{}).
		Save(chapters).
		Error

	return err
}
