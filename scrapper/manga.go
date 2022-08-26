package scrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func ScrapeMangaInfo(ctx web.Context, c models.Connector, opts *ScrapeOptions) (models.Manga, error) {
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
