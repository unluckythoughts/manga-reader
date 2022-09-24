package service

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) AddFavorite(ctx web.Context, link string) (models.MangaFavorite, error) {
	manga, err := s.db.GetManga(ctx, link)
	if err != nil {
		conn, err := connector.NewMangaConnector(ctx, link)
		if err != nil {
			return models.MangaFavorite{}, err
		}

		manga, err = conn.GetMangaInfo(ctx, link)
		if err != nil {
			return models.MangaFavorite{}, err
		}
	}

	favorite := models.MangaFavorite{
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

func (s *Service) GetFavorites(ctx web.Context) ([]models.MangaFavorite, error) {
	return s.db.GetFavorites(ctx)
}

func (s *Service) AddNovelFavorite(ctx web.Context, link string) (models.NovelFavorite, error) {
	novel, err := s.db.GetNovel(ctx, link)
	if err != nil {
		conn, err := connector.NewNovelConnector(ctx, link)
		if err != nil {
			return models.NovelFavorite{}, err
		}

		novel, err = conn.GetNovelInfo(ctx, link)
		if err != nil {
			return models.NovelFavorite{}, err
		}
	}

	favorite := models.NovelFavorite{
		Novel: novel,
	}

	return favorite, s.db.CreateNovelFavorite(ctx, &favorite)
}

func (s *Service) DelNovelFavorite(ctx web.Context, id int) error {
	favorite, err := s.db.FindNovelFavorite(ctx, id)
	if err != nil {
		return err
	}

	return s.db.DelNovelFavorite(ctx, favorite)
}

func (s *Service) GetNovelFavorites(ctx web.Context) ([]models.NovelFavorite, error) {
	return s.db.GetNovelFavorites(ctx)
}
