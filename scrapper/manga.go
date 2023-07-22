package scrapper

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
	"go.uber.org/zap"
)

func ScrapeMangaInfo(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) (models.Manga, error) {
	manga := models.Manga{}
	if err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		var err error
		manga, err = GetMangaFromInfoSelectors(h.DOM, c.Info)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
			return
		}
		manga.Source = c.Source
		manga.URL = opts.URL

		if c.Info.ChapterListLastPage != "" {
			manga.Chapters, err = _scrapeChaptersParallel(ctx, c, opts)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter info from %s", c.Source.Domain)
				return
			}
		} else {
			h.DOM.Find(c.Info.ChapterContainer).Each(func(i int, s *goquery.Selection) {
				chapter, err := GetChapterFromInfoSelectors(s, c.Info)
				if err != nil {
					ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter info from %s", c.Source.Domain)
					return
				}

				manga.Chapters = append(manga.Chapters, chapter)
			})
		}
	}); err != nil {
		return manga, err
	}

	return manga, nil
}

func _scrapeChaptersInPage(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) ([]models.MangaChapter, error) {
	chapters := []models.MangaChapter{}
	if err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		url := opts.URL
		h.DOM.Find(c.Info.ChapterContainer).Each(func(i int, s *goquery.Selection) {
			chapter, err := GetChapterFromInfoSelectors(s, c.Info)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting chapter info from %s", url)
				return
			}

			chapters = append(chapters, chapter)
		})
	}); err != nil {
		return chapters, err
	}

	return chapters, nil
}

func _scrapeChaptersParallel(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) ([]models.MangaChapter, error) {
	chapters := []models.MangaChapter{}

	var count = 0
	if strings.HasPrefix(c.Info.ChapterListLastPage, "DEFAULT::") {
		count = utils.GetInt(strings.Split(c.Info.ChapterListLastPage, "DEFAULT::")[1])
	} else {
		lastPage := _scrapLastPage(ctx, c.Info.ChapterListLastPage, opts)
		count = utils.GetInt(lastPage)
		if count == 0 {
			count = 1
		}
	}

	workerFn := func(page int64, out chan<- []models.MangaChapter) {
		params := strings.Replace(c.Info.ChapterListPageParam, PAGE_ID, strconv.Itoa(int(page)), 1)
		url := opts.URL + params
		newOpts := opts.Clone()
		newOpts.URL = url
		chapters, err := _scrapeChaptersInPage(ctx, c, newOpts)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error while getting mangas from page %s", url)
			return
		}
		ctx.Logger().Debugf("got %d chapters from page %d", len(chapters), page)

		out <- chapters
	}

	workerCount := 3
	payloadChan, chaptersChan := utils.StartWorkers[int64, []models.MangaChapter](workerCount, workerFn)
	utils.SendPayloads[int64](utils.MakeRange[int64](1, int64(count), 1), payloadChan)
	chapters = utils.GetSlicedResults[models.MangaChapter](chaptersChan)

	return chapters, nil
}
