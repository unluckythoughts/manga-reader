package scrapper

import (
	"errors"

	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Scrapper) populateManga(resp *mangaInfoResponse) func(h *colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		resp.Manga.Title, resp.Error = getText(h.DOM.Find(s.src.MangaInfo.TitleSelector), s.src.MangaList.NextPageSelector)
		if resp.Error != nil {
			return
		}

		resp.Manga.ImageURL, resp.Error = getText(h.DOM.Find(s.src.MangaInfo.ImageURLSelector), s.src.MangaList.NextPageSelector)
		if resp.Error != nil {
			return
		}

		resp.Manga.Synopsis, resp.Error = getText(h.DOM.Find(s.src.MangaInfo.SynopsisSelector), s.src.MangaList.NextPageSelector)
		if resp.Error != nil {
			return
		}

		var nums []string
		nums, resp.Error = getTextListForSelector(h, s.src.MangaInfo.ChapterNumberSelector)
		if resp.Error != nil {
			return
		}

		var urls []string
		urls, resp.Error = getTextListForSelector(h, s.src.MangaInfo.ChapterURLSelector)
		if resp.Error != nil {
			return
		}

		var titles []string
		titles, resp.Error = getTextListForSelector(h, s.src.MangaInfo.ChapterTitleSelector)
		if resp.Error != nil {
			return
		}

		var dates []string
		dates, resp.Error = getTextListForSelector(h, s.src.MangaInfo.ChapterUploadDateSelector)
		if resp.Error != nil {
			return
		}

		if (len(titles) > 0 && len(urls) != len(titles)) ||
			(len(nums) > 0 && len(urls) != len(nums)) ||
			(len(dates) > 0 && len(urls) != len(dates)) {
			resp.Error = errors.New("title and url number on list mismatch")
			return
		}

		for i := 0; i < len(urls); i++ {
			chapter := models.Chapter{
				URL: urls[i],
			}

			if len(titles) > 0 {
				chapter.Title = titles[i]
			}
			if len(nums) > 0 {
				chapter.Title = nums[i]
			}
			if len(dates) > 0 {
				chapter.Title = dates[i]
			}

			resp.Manga.Chapters = append(resp.Manga.Chapters, chapter)
		}
	}
}

type mangaInfoResponse struct {
	Manga    models.Manga
	Chapters []models.Chapter
	Error    error
}

func (s *Scrapper) GetMangaInfo(url string) (models.Manga, error) {
	c := s.getColly()

	if isInternalLink(url) {
		url = s.src.MangaList.URL + url
	}

	resp := mangaInfoResponse{
		Manga: models.Manga{
			URL: url,
		},
	}
	c.OnHTML("body", s.populateManga(&resp))

	err := c.Visit(url)
	if err != nil {
		return resp.Manga, err
	}

	return resp.Manga, resp.Error
}
