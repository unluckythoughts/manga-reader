package scrapper

import (
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func ScrapeChapterPages(ctx web.Context, c models.Connector, opts *ScrapeOptions) (models.Pages, error) {
	pages := models.Pages{}
	GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		imageURLs, err := GetImagesListForSelector(h.DOM, c.Chapter.ImageUrl, false)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter images from %s", c.Source.Domain)
			return
		}

		pages.URLs = imageURLs
	})

	return pages, nil
}
