package service

import (
	"net/http"
	"regexp"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
)

func (s *Service) GetMangaSourceList(ctx web.Context) ([]models.Source, error) {
	return s.db.GetSources(ctx)
}

func (s *Service) UpdateAllSourceMangas(ctx web.Context) error {
	return s.w.UpdateAllSourceMangas(ctx)
}

func (s *Service) GetSourceMangaList(ctx web.Context, id uint, force bool, fix bool) ([]models.Manga, error) {
	if !force {
		return s.db.GetSourceMangas(ctx, id)
	}

	src, err := s.db.FindSource(ctx, id)
	if err != nil {
		return []models.Manga{}, err
	}

	conn, err := connector.GetMangaConnector(ctx, src.Domain)
	if err != nil {
		return []models.Manga{}, err
	}
	mangas, err := conn.GetMangaList(ctx)

	if err == nil {
		s.w.UpdateSourceMangas(ctx, src, mangas, fix)
	}

	if len(mangas) > 400 {
		return mangas[:400], err
	}
	return mangas, err
}

func (s *Service) SearchSourceManga(ctx web.Context, query string) ([]models.Manga, error) {
	sitePattern := regexp.MustCompile(` ?site:([^\s]+) ?`)
	if !sitePattern.MatchString(query) {
		mangas, err := s.db.SearchMangasByTitle(ctx, query)
		return mangas, err
	}

	site := sitePattern.FindStringSubmatch(query)[1]
	query = sitePattern.ReplaceAllString(query, "")
	source, err := s.db.SearchSourceByDomain(ctx, site)
	if err != nil {
		if utils.GetInt(site) > 0 {
			sourceID := utils.GetInt(site)
			source, err = s.db.FindSource(ctx, uint(sourceID))
			if err != nil {
				return []models.Manga{}, err
			}
		}
	}

	mangas, err := s.db.SearchSourceMangasByTitle(ctx, source, query)
	return mangas, err

}

func (s *Service) GetSourceManga(ctx web.Context, id uint, force bool) (models.Manga, error) {
	dbmanga, err := s.db.GetManga(ctx, id)
	if !force {
		return dbmanga, err
	}
	source := dbmanga.Source

	conn, err := connector.GetMangaConnector(ctx, source.Domain)
	if err != nil {
		return models.Manga{}, err
	}

	manga, err := conn.GetMangaInfo(ctx, dbmanga.URL)
	if err != nil {
		if err.Error() == http.StatusText(http.StatusNotFound) {
			mangas, err := conn.GetMangaList(ctx)
			if err != nil {
				return manga, err
			}

			err = s.w.UpdateSourceMangasSync(ctx, source, mangas, true)
			if err != nil {
				return manga, err
			}

			for _, m := range mangas {
				if m.Title == dbmanga.Title {
					dbmanga.URL = manga.URL
					dbmanga.Chapters = manga.Chapters

					return manga, nil
				}
			}
		}

		return manga, err
	}

	err = s.w.UpdateSourceMangaSync(ctx, conn.GetSource().Domain, &manga)
	return manga, err
}

func (s *Service) GetSourceMangaChapter(ctx web.Context, id uint, force bool) (models.MangaChapter, error) {
	pages := models.Pages{}
	chapter, err := s.db.GetChapter(ctx, id)
	if err != nil {
		return chapter, err
	}

	if !force && len(chapter.ImageURLs) > 0 {
		return chapter, nil
	}

	domain, err := s.db.GetSourceDomain(ctx, chapter.MangaID)
	if err != nil {
		return chapter, err
	}

	conn, err := connector.GetMangaConnector(ctx, domain)
	if err != nil {
		return chapter, err
	}

	chapterURL := "https://" + domain + chapter.URL
	pages, err = conn.GetChapterPages(ctx, chapterURL)
	if err != nil {
		return chapter, err
	}
	chapter.ImageURLs = pages.URLs

	return chapter, s.db.UpdateChapterPages(ctx, chapter.ID, pages.URLs)
}

func (s *Service) GetNovelSourceList(ctx web.Context) ([]models.NovelSource, error) {
	return s.db.GetNovelSources(ctx)
}

func (s *Service) GetSourceNovelList(ctx web.Context, id uint, force bool) ([]models.Novel, error) {
	if !force {
		return s.db.GetSourceNovels(ctx, id)
	}

	src, err := s.db.FindSource(ctx, id)
	if err != nil {
		return []models.Novel{}, err
	}

	conn, err := connector.NewNovelConnector(ctx, src.Domain)
	if err != nil {
		return []models.Novel{}, err
	}
	mangas, err := conn.GetNovelList(ctx)

	if err == nil {
		s.w.UpdateSourceNovels(ctx, src.Domain, mangas)
	}

	if len(mangas) > 200 {
		return mangas[:200], err
	}
	return mangas, err
}

func (s *Service) SearchSourceNovel(ctx web.Context, query string) ([]models.Novel, error) {
	mangas, err := s.db.SearchNovelsByTitle(ctx, query)

	return mangas, err
}

func (s *Service) GetSourceNovel(ctx web.Context, id uint, force bool) (models.Novel, error) {
	manga, err := s.db.GetNovel(ctx, id)
	if !force {
		return manga, err
	}

	conn, err := connector.NewNovelConnector(ctx, manga.URL)
	if err != nil {
		return models.Novel{}, err
	}
	manga, err = conn.GetNovelInfo(ctx, manga.URL)

	if err == nil {
		s.w.UpdateSourceNovel(ctx, conn.GetSource().Domain, manga)
	}
	return manga, err
}

func (s *Service) GetSourceNovelChapter(ctx web.Context, id uint) ([]string, error) {
	chapter, err := s.db.GetNovelChapter(ctx, id)
	if err != nil {
		return []string{}, err
	}

	domain, err := s.db.GetSourceDomain(ctx, chapter.NovelID)
	if err != nil {
		return []string{}, err
	}

	chapterURL := "https://" + domain + chapter.URL

	conn, err := connector.NewNovelConnector(ctx, chapterURL)
	if err != nil {
		return []string{}, err
	}

	return conn.GetNovelChapter(ctx, chapterURL)
}
