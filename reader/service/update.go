package service

import (
	"time"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) UpdateMangaFavorite(ctx web.Context, favoriteID uint) (models.MangaFavorite, error) {
	favorite, err := s.db.FindFavorite(ctx, favoriteID)
	if err != nil {
		return models.MangaFavorite{}, err
	}

	conn, err := connector.GetMangaConnector(ctx, favorite.Manga.Source.Domain)
	if err != nil {
		return models.MangaFavorite{}, err
	}

	manga, err := conn.GetMangaInfo(ctx, favorite.Manga.URL)
	if err != nil {
		return models.MangaFavorite{}, err
	}

	manga.SourceID = favorite.Manga.SourceID
	err = s.db.UpdateMangas(ctx, &[]models.Manga{manga})
	if err != nil {
		ctx.Logger().
			With("error", err).
			Warnf("error updating mangas for %s", favorite.Manga.Title)
	}

	for i := range manga.Chapters {
		manga.Chapters[i].MangaID = favorite.MangaID
		if manga.Chapters[i].UploadDate == "" {
			manga.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
		}
	}

	// will return updated chapters
	err = s.db.UpdateChapters(ctx, &manga.Chapters)
	if err != nil {
		return models.MangaFavorite{}, err
	}

	favorite.Manga = manga

	return favorite, nil
}

func (s *Service) UpdateMangaFavoriteProgress(ctx web.Context, favoriteID uint, progress models.StrFloatList) error {
	return s.db.UpdateFavoriteProgress(ctx, favoriteID, progress)
}

func (s *Service) UpdateAllMangaFavorite(ctx web.Context) error {
	return s.w.UpdateFavorites(ctx)
}

func (s *Service) UpdateNovelFavorite(ctx web.Context, favoriteID uint) (models.NovelFavorite, error) {
	favorite, err := s.db.FindNovelFavorite(ctx, favoriteID)
	if err != nil {
		return models.NovelFavorite{}, err
	}

	conn, err := connector.NewNovelConnector(ctx, favorite.Novel.URL)
	if err != nil {
		return models.NovelFavorite{}, err
	}

	novel, err := conn.GetNovelInfo(ctx, favorite.Novel.URL)
	if err != nil {
		return models.NovelFavorite{}, err
	}

	novel.SourceID = favorite.Novel.SourceID
	err = s.db.UpdateNovels(ctx, &[]models.Novel{novel})
	if err != nil {
		ctx.Logger().
			With("error", err).
			Warnf("error updating novels for %s", favorite.Novel.Title)
	}

	for i := range novel.Chapters {
		novel.Chapters[i].NovelID = favorite.NovelID
		if novel.Chapters[i].UploadDate == "" {
			novel.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
		}
	}

	// will return updated chapters
	err = s.db.UpdateNovelChapters(ctx, &novel.Chapters)
	if err != nil {
		return models.NovelFavorite{}, err
	}

	favorite.Novel = novel

	return favorite, nil
}

func (s *Service) UpdateNovelFavoriteProgress(ctx web.Context, favoriteID uint, progress models.StrFloatList) error {
	return s.db.UpdateNovelFavoriteProgress(ctx, favoriteID, progress)
}

func (s *Service) UpdateAllNovelFavorite(ctx web.Context) error {
	return s.w.UpdateNovelFavorites(ctx)
}
