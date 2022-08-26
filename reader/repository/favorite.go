package repository

import (
	"context"

	"github.com/unluckythoughts/manga-reader/models"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateFavorite(ctx context.Context, f *models.Favorite) error {
	return r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true, UpdateAll: false}).
		Clauses(clause.Returning{}).
		Create(f).
		Error
}

func (r *Repository) GetFavorites(ctx context.Context) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := r.db.Model(&models.Favorite{}).Preload("Manga.Chapters").Find(&favorites).Error
	return favorites, err
}

func (r *Repository) FindFavorite(ctx context.Context, id int) (models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).
		Preload("Manga.Chapters").
		Preload("Manga.Source").
		First(&favorite, id).Error
	return favorite, err
}

func (r *Repository) UpdateFavoriteProgress(ctx context.Context, favoriteID int, progress models.StrFloatList) error {
	favorite := &models.Favorite{ID: favoriteID, Progress: progress}
	return r.db.
		WithContext(ctx).
		Updates(favorite).
		Error
}

func (r *Repository) DelFavorite(ctx context.Context, favorite models.Favorite) error {
	return r.db.
		WithContext(ctx).
		Delete(&models.Favorite{}, favorite.ID).
		Error
}
