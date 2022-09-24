package repository

import (
	"context"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateFavorite(ctx context.Context, f *models.MangaFavorite) error {
	return r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true, UpdateAll: false}).
		Clauses(clause.Returning{}).
		Create(f).
		Error
}

func (r *Repository) GetFavorites(ctx context.Context) ([]models.MangaFavorite, error) {
	var favorites []models.MangaFavorite
	err := r.db.
		Model(&models.MangaFavorite{}).
		Preload("Manga.Chapters").
		Preload("Manga.Source").
		Find(&favorites).Error
	return favorites, err
}

func (r *Repository) FindFavorite(ctx context.Context, id int) (models.MangaFavorite, error) {
	var favorite models.MangaFavorite
	err := r.db.WithContext(ctx).
		Preload("Manga.Chapters").
		Preload("Manga.Source").
		First(&favorite, id).Error
	return favorite, err
}

func (r *Repository) UpdateFavoriteProgress(ctx context.Context, favoriteID int, progress models.StrFloatList) error {
	favorite := &models.MangaFavorite{ID: favoriteID, Progress: progress}
	return r.db.
		WithContext(ctx).
		Updates(favorite).
		Error
}

func (r *Repository) DelFavorite(ctx context.Context, favorite models.MangaFavorite) error {
	return r.db.
		WithContext(ctx).
		Delete(&models.MangaFavorite{}, favorite.ID).
		Error
}
