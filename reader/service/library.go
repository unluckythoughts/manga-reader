package service

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) AddFavorite(ctx web.Context, link string) (models.Favorite, error) {
	manga, err := s.db.GetManga(ctx, link)
	if err != nil {
		conn, err := connector.NewMangaConnector(ctx, link)
		if err != nil {
			return models.Favorite{}, err
		}

		manga, err = conn.GetMangaInfo(ctx, link)
		if err != nil {
			return models.Favorite{}, err
		}
	}

	favorite := models.Favorite{
		Manga: manga,
	}

	return favorite, s.db.CreateFavorite(ctx, &favorite)
}

func (s *Service) DelFavorite(ctx web.Context, id int) error {
	favorite, err := s.db.FindFavorite(ctx, id)
	if err != nil {
		return err
	}

	return s.db.DelFavorite(ctx, favorite)
}

func (s *Service) GetFavorites(ctx web.Context) ([]models.Favorite, error) {
	return s.db.GetFavorites(ctx)
}
