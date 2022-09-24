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

func _scrapNovelChapterList(ctx web.Context, c models.NovelConnector, h *colly.HTMLElement) ([]models.NovelChapter, string) {
	chapters := []models.NovelChapter{}

	h.DOM.Find(c.NovelSelectors.Info.ChapterContainer).Each(func(i int, s *goquery.Selection) {
		chapter, err := GetNovelChapterFromInfoSelectors(s, c.NovelSelectors.Info)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter info from %s", c.Domain)
			return
		}

		if strings.HasPrefix(chapter.URL, "/") {
			chapter.URL, _ = url.JoinPath(c.BaseURL, chapter.URL)
		}

		chapters = append(chapters, chapter)
	})

	if c.NovelSelectors.Info.ChapterListNextPage == "" {
		return chapters, ""
	}

	nextPageURL, err := GetTextForSelector(h.DOM, c.NovelSelectors.Info.ChapterListNextPage)
	if err != nil {
		ctx.Logger().With(zap.Error(err)).Debug("error getting next chapters page url")
		return chapters, ""
	}

	if strings.HasPrefix(nextPageURL, "/") {
		nextPageURL, _ = url.JoinPath(c.BaseURL, nextPageURL)
	}

	return chapters, nextPageURL
}

func ScrapeNovelInfo(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions) (models.Novel, error) {
	novel := models.Novel{}
	err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		var err error
		novel, err = GetNovelFromInfoSelectors(h.DOM, c.NovelSelectors.Info)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting novel info from %s", c.Domain)
			return
		}
		novel.URL = opts.URL

		if c.NovelSelectors.Info.ChapterListURL == "" {
			novel.Chapters, _ = _scrapNovelChapterList(ctx, c, h)
			return
		}

		nextPageURL := c.NovelSelectors.Info.ChapterListURL
		for nextPageURL != "" {
			chapters := []models.NovelChapter{}
			if strings.HasPrefix(nextPageURL, "/") {
				opts.URL, _ = url.JoinPath(c.BaseURL, nextPageURL)
			} else if !strings.HasPrefix(nextPageURL, "http") {
				opts.URL, _ = url.JoinPath(opts.URL, nextPageURL)
			} else {
				opts.URL = nextPageURL
			}

			if err := GetPageForScrapping(ctx, opts, func(hh *colly.HTMLElement) {
				chapters, nextPageURL = _scrapNovelChapterList(ctx, c, hh)
			}); err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting chapters of novel from %s", c.Domain)
				return
			}

			novel.Chapters = append(novel.Chapters, chapters...)
		}
	})

	return novel, err
}
