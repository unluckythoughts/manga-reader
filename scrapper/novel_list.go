package scrapper

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
	"go.uber.org/zap"
)

func _scrapNovelsLastPage(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions) string {
	lastpage := ""
	err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		var selErr error
		lastpage, selErr = GetTextForSelector(h.DOM, c.NovelSelectors.List.LastPage)
		ctx.Logger().With(zap.Error(selErr)).Debugf("error getting last page text from %s", c.Domain)
	})

	if err != nil {
		ctx.Logger().With(zap.Error(err)).Debugf("error getting last page info from %s", c.Domain)
	}

	return lastpage
}

func _scrapNovelsInPage(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions) ([]models.Novel, string, error) {
	novels := []models.Novel{}
	nexPageURL := ""
	err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		h.DOM.Find(c.NovelSelectors.List.NovelContainer).Each(func(i int, s *goquery.Selection) {
			novel, err := GetNovelFromListSelectors(s, c.NovelSelectors.List)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting novel info from %s", c.Domain)
				return
			}

			if strings.HasPrefix(novel.URL, "/") {
				novel.URL, _ = url.JoinPath(c.BaseURL, novel.URL)
			}
			if strings.HasPrefix(novel.ImageURL, "/") {
				novel.ImageURL, _ = url.JoinPath(c.BaseURL, novel.ImageURL)
			}

			if novel.Title != "" {
				novels = append(novels, novel)
			}
		})

		if nextPageElement, ok := GetElementForSelector(h.DOM, c.NovelSelectors.List.NextPage); ok {
			if link, ok := nextPageElement.Attr("href"); ok {
				if strings.HasPrefix(link, c.NovelListPath) || strings.HasPrefix(link, "/") {
					nexPageURL, _ = url.JoinPath(c.BaseURL, link)
				} else {
					nexPageURL, _ = url.JoinPath(c.BaseURL, c.NovelListPath, link)
				}
			}
		}
	})

	if err != nil {
		ctx.Logger().With(zap.Error(err)).Debugf("error getting novel info from %s", c.Domain)
	}

	return novels, nexPageURL, err
}

func ScrapeNovels(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions) ([]models.Novel, error) {
	novels := []models.Novel{}
	// url := c.BaseURL + c.NovelListPath
	for opts.URL != "" {
		pageNovels, nextPage, err := _scrapNovelsInPage(ctx, c, opts)
		if err != nil {
			return novels, err
		}

		novels = append(novels, pageNovels...)
		opts.URL = nextPage
	}

	return novels, nil
}

func ScrapeNovelsParallel(ctx web.Context, c models.NovelConnector, opts *ScrapeOptions, threadCount int) ([]models.Novel, error) {
	novels := []models.Novel{}

	lastPage := _scrapNovelsLastPage(ctx, c, opts)
	lastPage = strings.Split(lastPage, " ")[len(strings.Split(lastPage, " "))-1]
	count := utils.GetInt(lastPage)

	workerFn := func(page int64, out chan<- []models.Novel) {
		params := strings.Replace(c.List.PageParam, MANGA_LIST_PAGE_ID, strconv.Itoa(int(page)), 1)
		url := c.BaseURL + c.NovelListPath + params
		newOpts := opts.Clone()
		newOpts.URL = url
		novels, _, err := _scrapNovelsInPage(ctx, c, newOpts)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error while getting novels from page %s", url)
			return
		}

		out <- novels
	}

	payloadChan, novelsChan := utils.StartWorkers[int64, []models.Novel](threadCount, workerFn)
	utils.SendPayloads[int64](utils.MakeRange[int64](1, int64(count), 1), payloadChan)
	novels = utils.GetSlicedResults[models.Novel](novelsChan)

	return novels, nil
}
