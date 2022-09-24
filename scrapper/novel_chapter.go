package scrapper

import (
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func ScrapeNovelChapterText(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions) ([]string, error) {
	text := []string{}
	GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		data, err := GetAllTextForSelector(h.DOM, c.Chapter.Paragraph)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter images from %s", c.Source.Domain)
			return
		}

		text = data
	})

	return text, nil
}
