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
		First(&favorite, id).Error
	return favorite, err
}

func (r *Repository) FindFavoriteByMangaURL(ctx context.Context, mangaUrl string) (models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).
		Joins("Manga").
		Where("Manga.url = ?", mangaUrl).
		First(&favorite).Error
	return favorite, err
}

func (r *Repository) UpdateFavoriteProgress(ctx context.Context, favoriteID int, progress models.StrIntList) error {
	favorite := &models.Favorite{ID: favoriteID, Progress: progress}
	return r.db.
		WithContext(ctx).
		Updates(favorite).
		Error
}

func (r *Repository) DelFavorite(ctx context.Context, favorite models.Favorite) error {
	tx := r.db.WithContext(ctx).Begin()

	err := tx.Delete(&models.Favorite{}, favorite.ID).Error
	if err != nil {
		_ = tx.Rollback().Error
		return err
	}

	err = tx.Delete(&models.Manga{}, favorite.MangaID).Error
	if err != nil {
		_ = tx.Rollback().Error
		return err
	}

	err = tx.Where("manga_id = ?", favorite.MangaID).Delete(&models.Chapter{}).Error
	if err != nil {
		_ = tx.Rollback().Error
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) UpdateFavoriteChapters(ctx context.Context, chapters []models.Chapter) error {
	return r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			UpdateAll: true,
		}).
		Save(&chapters).
		Error
}
