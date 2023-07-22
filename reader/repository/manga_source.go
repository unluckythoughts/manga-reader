package repository

import (
	"context"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateSources(ctx context.Context, sources *[]models.Source) error {
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

func (r *Repository) SaveSource(ctx context.Context, source *models.Source) error {
	return r.db.
		WithContext(ctx).
		Save(source).
		Error
}

func (r *Repository) FindSource(ctx context.Context, id uint) (models.Source, error) {
	var source models.Source
	err := r.db.WithContext(ctx).
		First(&source, id).Error
	return source, err
}

// function to search the name by domain
func (r *Repository) SearchSourceByDomain(ctx context.Context, query string) (models.Source, error) {
	var source models.Source
	err := r.db.WithContext(ctx).
		Where("LOWER(domain) REGEXP ?", strings.ToLower(query)).
		First(&source).Error

	return source, err
}

func (r *Repository) FindSourceByDomain(ctx context.Context, domain string) (models.Source, error) {
	var source models.Source
	err := r.db.WithContext(ctx).
		Where("domain = ?", domain).
		First(&source).Error
	return source, err
}

func (r *Repository) GetSources(ctx context.Context) ([]models.Source, error) {
	var sources []models.Source
	err := r.db.
		Model(&models.Source{}).
		Find(&sources).
		Error

	return sources, err
}
