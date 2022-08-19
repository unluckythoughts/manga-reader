package service

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) AddFavorite(ctx web.Context, link string) error {
	conn, err := connector.New(ctx, link)
	if err != nil {
		return err
	}

	manga, err := conn.GetMangaInfo(ctx, link)
	if err != nil {
		return err
	}

	favorite := models.Favorite{
		Manga: manga,
	}

	return s.db.CreateFavorite(ctx, &favorite)
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
