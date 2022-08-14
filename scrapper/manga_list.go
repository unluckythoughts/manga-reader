package scrapper

import (
	"errors"

	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/manga-reader/models"
)

func (s *Scrapper) populateMangas(resp *mangaListResponse) func(h *colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		var titles []string
		titles, resp.Error = getTextListForSelector(h, s.src.MangaList.MangaTitleSelector)
		if resp.Error != nil {
			return
		}

		var urls []string
		urls, resp.Error = getTextListForSelector(h, s.src.MangaList.MangaURLSelector)
		if resp.Error != nil {
			return
		}

		var imageUrls []string
		imageUrls, resp.Error = getTextListForSelector(h, s.src.MangaList.MangaImageURLSelector)
		if resp.Error != nil {
			return
		}

		if (len(titles) > 0 && len(urls) != len(titles)) ||
			(len(imageUrls) > 0 && len(urls) != len(imageUrls)) {
			resp.Error = errors.New("title and url number on list mismatch")
			return
		}

		for i := 0; i < len(urls); i++ {

			manga := models.Manga{
				URL: urls[i],
			}

			if len(titles) > 0 {
				manga.Title = titles[i]
			}
			if len(imageUrls) > 0 {
				manga.ImageURL = imageUrls[i]
			}

			resp.Mangas = append(resp.Mangas, manga)
		}

		// checking if the next page exists
		resp.NextPage, resp.Error = getText(h.DOM.Find(s.src.MangaList.NextPageSelector), s.src.MangaList.NextPageSelector)
	}
}

type mangaListResponse struct {
	Mangas   []models.Manga
	NextPage string
	Error    error
}

func (s *Scrapper) GetMangaList() ([]models.Manga, error) {
	c := s.getColly()

	resp := mangaListResponse{
		NextPage: s.src.MangaList.URL,
	}
	c.OnHTML("body", s.populateMangas(&resp))

	for resp.NextPage != "" {
		if isInternalLink(resp.NextPage) {
			resp.NextPage = s.src.MangaList.URL + resp.NextPage
		}

		err := c.Visit(resp.NextPage)
		if err != nil {
			return resp.Mangas, err
		}

		if resp.Error != nil {
			return resp.Mangas, resp.Error
		}
	}

	return resp.Mangas, resp.Error
}
