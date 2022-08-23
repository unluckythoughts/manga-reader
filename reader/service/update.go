package service

import (
	"time"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

// TODO: reuse source manga and delete this
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

	chapters := append(favorite.Manga.Chapters, manga.Chapters...)
	favorite.Manga = manga
	favorite.Manga.Chapters = chapters

	return favorite, nil
}

func (s *Service) UpdateFavoriteProgress(ctx web.Context, favoriteID int, progress models.StrIntList) error {
	return s.db.UpdateFavoriteProgress(ctx, favoriteID, progress)
}

// TODO: move to worker and delete this
func (s *Service) UpdateAllFavorite(ctx web.Context) error {
	favorites, err := s.db.GetFavorites(ctx)
	if err != nil {
		return err
	}

	for _, favorite := range favorites {
		ctx.Logger().Debugf("Updating manga %s", favorite.Manga.Title)

		conn, err := connector.New(ctx, favorite.Manga.URL)
		if err != nil {
			return err
		}
		manga, err := conn.GetMangaInfo(ctx, favorite.Manga.URL)
		if err != nil {
			ctx.Logger().
				With("error", err).
				Warnf("error scrapping chapters for %s", favorite.Manga.Title)
			continue
		}

		for i := range manga.Chapters {
			manga.Chapters[i].MangaID = favorite.MangaID
			if manga.Chapters[i].UploadDate == "" {
				manga.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
			}
		}

		err = s.db.UpdateChapters(ctx, &manga.Chapters)
		if err != nil {
			ctx.Logger().
				With("error", err).
				Warnf("error updating chapters for %s", favorite.Manga.Title)
		}
	}

	return nil
}
