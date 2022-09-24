package repository

import (
	"context"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetSourceNovels(ctx context.Context, domain string) ([]models.Novel, error) {
	var novels []models.Novel
	err := r.db.WithContext(ctx).
		Joins("Source").
		Where("Source.domain = ?", domain).
		Preload("Source").
		Limit(200).
		Find(&novels).
		Error

	return novels, err
}

func (r *Repository) SearchNovelsByTitle(ctx context.Context, query string) ([]models.Novel, error) {
	var novels []models.Novel
	err := r.db.WithContext(ctx).
		Where("LOWER(title) REGEXP ?", strings.ToLower(query)).
		Preload("Source").
		Find(&novels).
		Limit(100).
		Error

	return novels, err
}

func (r *Repository) UpdateNovels(ctx context.Context, novels *[]models.Novel) error {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			DoUpdates: clause.AssignmentColumns([]string{"synopsis", "image_url", "slug", "other_id"}),
		}, clause.Returning{}).
		Save(novels).
		Error

	return err
}

func (r *Repository) GetNovel(ctx context.Context, url string) (models.Novel, error) {
	var novel models.Novel
	err := r.db.
		Where(&models.Novel{URL: url}).
		Preload("Source").
		Preload("Chapters").
		Find(&novel).
		Error

	return novel, err
}

func (r *Repository) UpdateNovelChapters(ctx context.Context, chapters *[]models.NovelChapter) error {
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
