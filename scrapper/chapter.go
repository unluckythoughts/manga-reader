package scrapper

import (
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func ScrapeChapterPages(ctx web.Context, sels models.ChapterInfoSelectors, rt http.RoundTripper) ([]string, error) {
	c := getColly(ctx, rt)

	var imageURLs []string
	var chapterError error
	c.OnHTML("body", func(h *colly.HTMLElement) {
		imageURLs, chapterError = getTextListForSelector(h, sels.PageSelector)
	})

	err := c.Visit(sels.URL)
	if err != nil {
		return imageURLs, err
	}

	return imageURLs, chapterError
}
