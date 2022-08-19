package service

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) UpdateFavorite(ctx web.Context, favoriteID int) (models.Manga, error) {
	favorite, err := s.db.FindFavorite(ctx, favoriteID)
	if err != nil {
		return models.Manga{}, err
	}

	manga, err := s.GetSourceManga(ctx, favorite.Manga.URL)
	if err != nil {
		return models.Manga{}, err
	}

	for i := range manga.Chapters {
		manga.Chapters[i].MangaID = favorite.MangaID
	}

	err = s.db.UpdateFavoriteChapters(ctx, manga.Chapters)
	if err != nil {
		return models.Manga{}, err
	}

	return manga, nil
}

func (s *Service) UpdateFavoriteProgress(ctx web.Context, favoriteID int, progress models.StrIntList) error {
	return s.db.UpdateFavoriteProgress(ctx, favoriteID, progress)
}

func (s *Service) UpdateAllFavorite(ctx web.Context) error {
	favorites, err := s.db.GetFavorites(ctx)
	if err != nil {
		return err
	}

	for _, favorite := range favorites {
		manga, err := s.GetSourceManga(ctx, favorite.Manga.URL)
		if err != nil {
			ctx.Logger().
				With("error", err).
				Warnf("error scrapping chapters for %s", favorite.Manga.Title)
			continue
		}

		for i := range manga.Chapters {
			manga.Chapters[i].MangaID = favorite.MangaID
		}

		err = s.db.UpdateFavoriteChapters(ctx, manga.Chapters)
		if err != nil {
			ctx.Logger().
				With("error", err).
				Warnf("error updating chapters for %s", favorite.Manga.Title)
		}
	}

	return nil
}
