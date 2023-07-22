package repository

import (
	"context"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateNovelFavorite(ctx context.Context, f *models.NovelFavorite) error {
	return r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true, UpdateAll: false}).
		Clauses(clause.Returning{}).
		Create(f).
		Error
}

func (r *Repository) GetNovelFavorites(ctx context.Context) ([]models.NovelFavorite, error) {
	var favorites []models.NovelFavorite
	err := r.db.
		Model(&models.NovelFavorite{}).
		Preload("Novel.Chapters").
		Preload("Novel.Source").
		Find(&favorites).Error
	return favorites, err
}

func (r *Repository) FindNovelFavorite(ctx context.Context, id uint) (models.NovelFavorite, error) {
	var favorite models.NovelFavorite
	err := r.db.WithContext(ctx).
		Preload("Novel.Chapters").
		Preload("Novel.Source").
		First(&favorite, id).Error
	return favorite, err
}

func (r *Repository) UpdateNovelFavoriteProgress(ctx context.Context, favoriteID uint, progress models.StrFloatList) error {
	favorite := &models.NovelFavorite{ID: favoriteID, Progress: progress}
	return r.db.
		WithContext(ctx).
		Updates(favorite).
		Error
}

func (r *Repository) DelNovelFavorite(ctx context.Context, favorite models.NovelFavorite) error {
	return r.db.
		WithContext(ctx).
		Delete(&models.NovelFavorite{}, favorite.ID).
		Error
}
