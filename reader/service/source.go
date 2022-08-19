package service

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) GetSourceMangaList(ctx web.Context, link string) ([]models.Manga, error) {
	conn, err := connector.New(ctx, link)
	if err != nil {
		return []models.Manga{}, err
	}

	return conn.GetMangaList(ctx)
}

func (s *Service) GetSourceManga(ctx web.Context, mangaURL string) (models.Manga, error) {
	conn, err := connector.New(ctx, mangaURL)
	if err != nil {
		return models.Manga{}, nil
	}

	return conn.GetMangaInfo(ctx, mangaURL)
}

func (s *Service) GetSourceMangaChapter(ctx web.Context, chapterURL string) ([]string, error) {
	conn, err := connector.New(ctx, chapterURL)
	if err != nil {
		return []string{}, nil
	}

	return conn.GetChapterPages(ctx, chapterURL)
}
