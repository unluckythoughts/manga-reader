package scrapper

import (
	"github.com/gocolly/colly/v2"
)

func (s *Scrapper) GetChapterImageURLs(url string) ([]string, error) {
	c := s.getColly()

	if isInternalLink(url) {
		url = s.src.MangaList.URL + url
	}

	var imageURLs []string
	var chapterError error
	c.OnHTML("body", func(h *colly.HTMLElement) {
		imageURLs, chapterError = getTextListForSelector(h, s.src.ChapterInfo.ImageURLsSelector)
	})

	err := c.Visit(url)
	if err != nil {
		return imageURLs, err
	}

	return imageURLs, chapterError
}
