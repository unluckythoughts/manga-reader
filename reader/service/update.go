package service

import (
	"time"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) UpdateFavorite(ctx web.Context, favoriteID int) (models.Favorite, error) {
	favorite, err := s.db.FindFavorite(ctx, favoriteID)
	if err != nil {
		return models.Favorite{}, err
	}

	conn, err := connector.New(ctx, favorite.Manga.URL)
	if err != nil {
		return models.Favorite{}, err
	}

	manga, err := conn.GetMangaInfo(ctx, favorite.Manga.URL)
	if err != nil {
		return models.Favorite{}, err
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
		return models.Favorite{}, err
	}

	favorite.Manga = manga

	return favorite, nil
}

func (s *Service) UpdateFavoriteProgress(ctx web.Context, favoriteID int, progress models.StrFloatList) error {
	return s.db.UpdateFavoriteProgress(ctx, favoriteID, progress)
}

func (s *Service) UpdateAllFavorite(ctx web.Context) error {
	return s.w.UpdateFavorites(ctx)
}
