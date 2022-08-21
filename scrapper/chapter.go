package scrapper

import (
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func ScrapeChapterPages(ctx web.Context, sels models.ChapterInfoSelectors, opts *ScrapeOptions) ([]string, error) {
	opts.SetDefaults()
	c := getColly(ctx, opts.RoundTripper)

	var imageURLs []string
	var chapterError error
	c.OnHTML(opts.InitialHtmlTag, func(h *colly.HTMLElement) {
		imageURLs, chapterError = getTextListForSelector(h, sels.PageSelector)
	})

	err := c.Request(opts.RequestMethod, sels.URL, opts.Body, nil, opts.Headers)
	if err != nil {
		return imageURLs, err
	}

	return imageURLs, chapterError
}
