package repository

import (
	"context"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetSourceNovels(ctx context.Context, id uint) ([]models.Novel, error) {
	var novels []models.Novel
	err := r.db.WithContext(ctx).
		Joins("Source").
		Where("source_id = ?", id).
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

func (r *Repository) GetNovel(ctx context.Context, id uint) (models.Novel, error) {
	var novel models.Novel
	err := r.db.
		Where(id).
		Preload("Source").
		Preload("Chapters").
		Find(&novel).
		Error

	return novel, err
}

func (r *Repository) GetNovelChapter(ctx context.Context, id uint) (models.NovelChapter, error) {
	var chapter models.NovelChapter
	err := r.db.
		Where(id).
		Find(&chapter).
		Error

	return chapter, err
}

func (r *Repository) GetNovelSourceDomain(ctx context.Context, novelID uint) (string, error) {
	var novel models.Novel
	err := r.db.
		Where(novelID).
		Find(&novel).
		Error

	return novel.Source.Domain, err
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
