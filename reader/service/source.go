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

	conn, err := connector.NewMangaConnector(ctx, domain)
	if err != nil {
		return []models.Manga{}, err
	}
	mangas, err := conn.GetMangaList(ctx)

	if err == nil {
		s.w.UpdateSourceMangas(ctx, domain, mangas)
	}

	if len(mangas) > 200 {
		return mangas[:200], err
	}
	return mangas, err
}

func (s *Service) SearchSourceManga(ctx web.Context, query string) ([]models.Manga, error) {
	mangas, err := s.db.SearchMangasByTitle(ctx, query)

	return mangas, err
}

func (s *Service) GetSourceManga(ctx web.Context, mangaURL string, force bool) (models.Manga, error) {
	if !force {
		return s.db.GetManga(ctx, mangaURL)
	}

	conn, err := connector.NewMangaConnector(ctx, mangaURL)
	if err != nil {
		return models.Manga{}, err
	}
	manga, err := conn.GetMangaInfo(ctx, mangaURL)

	if err == nil {
		s.w.UpdateSourceManga(ctx, conn.GetSource().Domain, manga)
	}
	return manga, err
}

func (s *Service) GetSourceMangaChapter(ctx web.Context, chapterURL string) (models.Pages, error) {
	conn, err := connector.NewMangaConnector(ctx, chapterURL)
	if err != nil {
		return models.Pages{}, err
	}

	return conn.GetChapterPages(ctx, chapterURL)
}
