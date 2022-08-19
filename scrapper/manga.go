package scrapper

import (
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func populateManga(ctx web.Context, sels models.MangaInfoSelectors, resp *mangaInfoResponse) func(h *colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		resp.Manga.Title, resp.Error = getTextForSelector(h, sels.TitleSelector)
		if resp.Error != nil {
			return
		}

		resp.Manga.ImageURL, resp.Error = getTextForSelector(h, sels.ImageURLSelector)
		if resp.Error != nil {
			return
		}

		resp.Manga.Synopsis, resp.Error = getTextForSelector(h, sels.SynopsisSelector)
		if resp.Error != nil {
			return
		}

		var nums []string
		nums, resp.Error = getTextListForSelector(h, sels.ChapterNumberSelector)
		if resp.Error != nil {
			return
		}

		var urls []string
		urls, resp.Error = getTextListForSelector(h, sels.ChapterURLSelector)
		if resp.Error != nil {
			return
		}

		var titles []string
		titles, resp.Error = getTextListForSelector(h, sels.ChapterTitleSelector)
		if resp.Error != nil {
			return
		}

		var dates []string
		dates, resp.Error = getTextListForSelector(h, sels.ChapterUploadDateSelector)
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
				chapter.Number = GetChapterNumber(nums[i])
			}
			if len(dates) > 0 {
				chapter.UploadDate = dates[i]
			}

			resp.Manga.Chapters = append(resp.Manga.Chapters, chapter)
		}
	}
}

// func (s *Scrapper) populateMangaFromAPI(resp *mangaInfoResponse) {
// 	c := web.NewClient(resp.Manga.URL)

// 	apiResp := s.src.APIData.MangaData.APIQueryData.Response
// 	status, err := c.GetResponse("/", apiResp)
// 	if err != nil {
// 		resp.Error = errors.Wrapf(err, "error while get data from %s", resp.Manga.URL)
// 		return
// 	}

// 	if status != 200 {
// 		resp.Error = errors.Errorf("unexpected status %d when get data from %s", status, resp.Manga.URL)
// 		return
// 	}

// 	resp.Manga = s.src.APIData.MangaData.GetManga(apiResp)
// }

type mangaInfoResponse struct {
	Manga models.Manga
	Error error
}

func ScrapeMangaInfo(ctx web.Context, sels models.MangaInfoSelectors, rt http.RoundTripper) (models.Manga, error) {
	resp := mangaInfoResponse{
		Manga: models.Manga{
			URL: sels.URL,
		},
	}
	c := getColly(ctx, rt)

	c.OnHTML("body", populateManga(ctx, sels, &resp))

	err := c.Visit(sels.URL)
	if err != nil {
		return resp.Manga, err
	}

	return resp.Manga, resp.Error

}
