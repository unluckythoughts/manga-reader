package scrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func ScrapeMangaInfo(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) (models.Manga, error) {
	manga := models.Manga{}
	if err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		var err error
		manga, err = GetMangaFromInfoSelectors(h.DOM, c.MangaSelectors.Info)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
			return
		}
		manga.URL = opts.URL

		h.DOM.Find(c.MangaSelectors.Info.ChapterContainer).Each(func(i int, s *goquery.Selection) {
			chapter, err := GetChapterFromInfoSelectors(s, c.MangaSelectors.Info)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter info from %s", c.Source.Domain)
				return
			}

			manga.Chapters = append(manga.Chapters, chapter)
		})
	}); err != nil {
		return manga, err
	}

	return manga, nil
}
