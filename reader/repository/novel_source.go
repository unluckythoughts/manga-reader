package repository

import (
	"context"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateNovelSources(ctx context.Context, sources *[]models.NovelSource) error {
	return r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: false, UpdateAll: true,
		}).
		Clauses(clause.Returning{}).
		Create(sources).
		Error
}

func (r *Repository) SaveNovelSource(ctx context.Context, source *models.NovelSource) error {
	return r.db.
		WithContext(ctx).
		Save(source).
		Error
}

func (r *Repository) FindNovelSource(ctx context.Context, id int) (models.NovelSource, error) {
	var source models.NovelSource
	err := r.db.WithContext(ctx).
		First(&source, id).Error
	return source, err
}

func (r *Repository) FindNovelSourceByDomain(ctx context.Context, domain string) (models.NovelSource, error) {
	var source models.NovelSource
	err := r.db.WithContext(ctx).
		Where("domain = ?", domain).
		First(&source).Error
	return source, err
}

func (r *Repository) GetNovelSources(ctx context.Context) ([]models.NovelSource, error) {
	var sources []models.NovelSource
	err := r.db.
		Model(&models.NovelSource{}).
		Find(&sources).
		Error

	return sources, err
}
