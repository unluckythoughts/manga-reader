package scrapper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func ScrapeMangas(ctx web.Context, c models.Connector, opts *ScrapeOptions) ([]models.Manga, error) {
	mangas := []models.Manga{}
	// url := c.BaseURL + c.MangaListPath
	for opts.URL != "" {
		err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
			h.DOM.Find(c.Selectors.List.MangaContainer).Each(func(i int, s *goquery.Selection) {
				manga, err := GetMangaFromListSelectors(s, c.Selectors.List)
				if err != nil {
					ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
					return
				}

				if manga.Title != "" {
					mangas = append(mangas, manga)
				}

			})

			opts.URL = ""
			if nextPageElement, ok := GetElementForSelector(h.DOM, c.Selectors.List.NextPage); ok {
				if link, ok := nextPageElement.Attr("href"); ok {
					if strings.HasPrefix(link, c.MangaListPath) {
						opts.URL = c.BaseURL + link
					} else {
						opts.URL = c.BaseURL + c.MangaListPath + link
					}
				}
			}
		})

		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
			break
		}
	}

	return mangas, nil
}

func NewScrapeMangaInfo(ctx web.Context, c models.Connector, opts *ScrapeOptions) (models.Manga, error) {
	manga := models.Manga{}
	GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		var err error
		manga, err = GetMangaFromInfoSelectors(h.DOM, c.Selectors.Info)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
			return
		}
		manga.URL = opts.URL

		h.DOM.Find(c.Selectors.Info.ChapterContainer).Each(func(i int, s *goquery.Selection) {
			chapter, err := GetChapterFromInfoSelectors(s, c.Selectors.Info)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter info from %s", c.Source.Domain)
				return
			}

			manga.Chapters = append(manga.Chapters, chapter)
		})
	})

	return manga, nil
}

func NewScrapeChapterPages(ctx web.Context, c models.Connector, opts *ScrapeOptions) (models.Pages, error) {
	pages := models.Pages{}
	GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		imageURLs, err := GetImagesListForSelector(h.DOM, c.Chapter.ImageUrl, true)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter images from %s", c.Source.Domain)
			return
		}

		pages.URLs = imageURLs
	})

	return pages, nil
}
