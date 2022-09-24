package scrapper

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func _scrapNovelsInPage(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions) ([]models.Novel, string, error) {
	novels := []models.Novel{}
	nexPageURL := ""
	err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		h.DOM.Find(c.NovelSelectors.List.NovelContainer).Each(func(i int, s *goquery.Selection) {
			novel, err := GetNovelFromListSelectors(s, c.NovelSelectors.List)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting novel info from %s", c.Source.Domain)
				return
			}

			if strings.HasPrefix(novel.URL, "/") {
				novel.URL, _ = url.JoinPath(c.BaseURL, novel.URL)
			}
			if strings.HasPrefix(novel.ImageURL, "/") {
				novel.ImageURL, _ = url.JoinPath(c.BaseURL, novel.ImageURL)
			}

			if novel.Title != "" {
				novels = append(novels, novel)
			}
		})

		if nextPageElement, ok := GetElementForSelector(h.DOM, c.NovelSelectors.List.NextPage); ok {
			if link, ok := nextPageElement.Attr("href"); ok {
				if strings.HasPrefix(link, c.NovelListPath) || strings.HasPrefix(link, "/") {
					nexPageURL, _ = url.JoinPath(c.BaseURL, link)
				} else {
					nexPageURL, _ = url.JoinPath(c.BaseURL, c.NovelListPath, link)
				}
			}
		}
	})

	if err != nil {
		ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
	}

	return novels, nexPageURL, err
}

func ScrapeNovels(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions) ([]models.Novel, error) {
	mangas := []models.Novel{}
	// url := c.BaseURL + c.NovelListPath
	for opts.URL != "" {
		pageNovels, nextPage, err := _scrapNovelsInPage(ctx, c, opts)
		if err != nil {
			return mangas, err
		}

		mangas = append(mangas, pageNovels...)
		opts.URL = nextPage
	}

	return mangas, nil
}
