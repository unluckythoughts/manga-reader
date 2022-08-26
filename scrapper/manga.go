package scrapper

import (
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
	"go.uber.org/zap"
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
		nums, resp.Error = getTextListForSelector(h, sels.ChapterNumberSelector, false)
		if resp.Error != nil {
			return
		}

		var urls []string
		urls, resp.Error = getTextListForSelector(h, sels.ChapterURLSelector, false)
		if resp.Error != nil {
			return
		}

		var titles []string
		titles, resp.Error = getTextListForSelector(h, sels.ChapterTitleSelector, false)
		if resp.Error != nil {
			return
		}

		var dates []string
		dates, resp.Error = getTextListForSelector(h, sels.ChapterUploadDateSelector, false)
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
				var err error
				chapter.UploadDate, err = utils.ParseDate(dates[i], sels.ChapterUploadDateFormat)
				if err != nil {
					ctx.Logger().With(zap.Error(err)).Debug("could not parse upload date %s", dates[i])
				}
			}

			resp.Manga.Chapters = append(resp.Manga.Chapters, chapter)
		}
	}
}

type mangaInfoResponse struct {
	Manga models.Manga
	Error error
}

func ScrapeMangaInfo(ctx web.Context, sels models.MangaInfoSelectors, opts *ScrapeOptions) (models.Manga, error) {
	opts.SetDefaults()
	resp := mangaInfoResponse{
		Manga: models.Manga{
			URL: sels.URL,
		},
	}
	c := getColly(ctx, opts.RoundTripper)

	c.OnHTML(opts.InitialHtmlTag, populateManga(ctx, sels, &resp))

	err := c.Request(opts.RequestMethod, sels.URL, opts.Body, nil, opts.Headers)
	if err != nil {
		return resp.Manga, err
	}

	return resp.Manga, resp.Error

}
