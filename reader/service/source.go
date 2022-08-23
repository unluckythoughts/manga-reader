package service

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Service) GetSourceList(ctx web.Context) ([]models.Source, error) {
	return s.db.GetSources(ctx)
}

func (s *Service) GetSourceMangaList(ctx web.Context, domain string, force bool) ([]models.Manga, error) {
	if !force {
		return s.db.GetSourceMangas(ctx, domain)
	}

	conn, err := connector.New(ctx, domain)
	if err != nil {
		return []models.Manga{}, err
	}
	mangas, err := conn.GetMangaList(ctx)

	s.w.UpdateSourceMangas(ctx, domain, mangas)
	return mangas, err
}

func (s *Service) GetSourceManga(ctx web.Context, mangaURL string, force bool) (models.Manga, error) {
	if !force {
		return s.db.GetManga(ctx, mangaURL)
	}

	conn, err := connector.New(ctx, mangaURL)
	if err != nil {
		return models.Manga{}, nil
	}
	manga, err := conn.GetMangaInfo(ctx, mangaURL)

	s.w.UpdateSourceManga(ctx, conn.GetSource().Domain, manga)
	return manga, err
}

func (s *Service) GetSourceMangaChapter(ctx web.Context, chapterURL string) ([]string, error) {
	conn, err := connector.New(ctx, chapterURL)
	if err != nil {
		return []string{}, nil
	}

	return conn.GetChapterPages(ctx, chapterURL)
}
