package service

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) AddMangaFavorite(ctx web.Context, id uint) (models.MangaFavorite, error) {
	manga, err := s.db.GetManga(ctx, id)
	if err != nil {
		conn, err := connector.GetMangaConnector(ctx, manga.Source.Domain)
		if err != nil {
			return models.MangaFavorite{}, err
		}

		manga, err = conn.GetMangaInfo(ctx, manga.URL)
		if err != nil {
			return models.MangaFavorite{}, err
		}
	}

	favorite := models.MangaFavorite{
		Manga: manga,
	}

	return favorite, s.db.CreateFavorite(ctx, &favorite)
}

func (s *Service) DelMangaFavorite(ctx web.Context, id uint) error {
	favorite, err := s.db.FindFavorite(ctx, id)
	if err != nil {
		return err
	}

	return s.db.DelFavorite(ctx, favorite)
}

func (s *Service) GetMangaFavorites(ctx web.Context) ([]models.MangaFavorite, error) {
	return s.db.GetFavorites(ctx)
}

func (s *Service) AddNovelFavorite(ctx web.Context, id uint) (models.NovelFavorite, error) {
	novel, err := s.db.GetNovel(ctx, id)
	if err != nil {
		conn, err := connector.NewNovelConnector(ctx, novel.URL)
		if err != nil {
			return models.NovelFavorite{}, err
		}

		novel, err = conn.GetNovelInfo(ctx, novel.URL)
		if err != nil {
			return models.NovelFavorite{}, err
		}
	}

	favorite := models.NovelFavorite{
		Novel: novel,
	}

	return favorite, s.db.CreateNovelFavorite(ctx, &favorite)
}

func (s *Service) DelNovelFavorite(ctx web.Context, id uint) error {
	favorite, err := s.db.FindNovelFavorite(ctx, id)
	if err != nil {
		return err
	}

	return s.db.DelNovelFavorite(ctx, favorite)
}

func (s *Service) GetNovelFavorites(ctx web.Context) ([]models.NovelFavorite, error) {
	return s.db.GetNovelFavorites(ctx)
}
