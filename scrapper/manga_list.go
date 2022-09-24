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

const (
	MANGA_LIST_PAGE_ID = "::pageId::"
)

func _scrapLastPage(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) string {
	lastPageURL := ""
	err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		if lastPageElement, ok := GetElementForSelector(h.DOM, c.MangaSelectors.List.LastPage); ok {
			if link, ok := lastPageElement.Attr("href"); ok {
				lastPageURL = link
			}
		}
	})

	if err != nil {
		ctx.Logger().With(zap.Error(err)).Debugf("error getting last page info from %s", c.Source.Domain)
	}

	return lastPageURL
}

func _scrapMangasInPage(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) ([]models.Manga, string, error) {
	mangas := []models.Manga{}
	nexPageURL := ""
	err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
		h.DOM.Find(c.MangaSelectors.List.MangaContainer).Each(func(i int, s *goquery.Selection) {
			manga, err := GetMangaFromListSelectors(s, c.MangaSelectors.List)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
				return
			}

			if manga.Title != "" {
				mangas = append(mangas, manga)
			}

		})

		if nextPageElement, ok := GetElementForSelector(h.DOM, c.MangaSelectors.List.NextPage); ok {
			if link, ok := nextPageElement.Attr("href"); ok {
				if strings.HasPrefix(link, c.MangaListPath) {
					nexPageURL = c.BaseURL + link
				} else {
					nexPageURL = c.BaseURL + c.MangaListPath + link
				}
			}
		}
	})

	if err != nil {
		ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
	}

	return mangas, nexPageURL, err
}

func ScrapeMangas(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) ([]models.Manga, error) {
	mangas := []models.Manga{}
	// url := c.BaseURL + c.MangaListPath
	for opts.URL != "" {
		pageMangas, nextPage, err := _scrapMangasInPage(ctx, c, opts)
		if err != nil {
			return mangas, err
		}

		mangas = append(mangas, pageMangas...)
		opts.URL = nextPage
	}

	return mangas, nil
}

func ScrapeMangasParallel(ctx web.Context, c models.MangaConnector, opts *ScrapeOptions) ([]models.Manga, error) {
	mangas := []models.Manga{}

	lastPage := _scrapLastPage(ctx, c, opts)
	count := utils.GetInt(lastPage)

	workerFn := func(page int64, out chan<- []models.Manga) {
		params := strings.Replace(c.List.PageParam, MANGA_LIST_PAGE_ID, strconv.Itoa(int(page)), 1)
		url := opts.URL + params
		newOpts := opts.Clone()
		newOpts.URL = url
		mangas, _, err := _scrapMangasInPage(ctx, c, newOpts)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error while getting mangas from page %s", url)
			return
		}

		out <- mangas
	}

	workerCount := 10
	payloadChan, mangasChan := utils.StartWorkers[int64, []models.Manga](workerCount, workerFn)
	utils.SendPayloads[int64](utils.MakeRange[int64](1, int64(count), 1), payloadChan)
	mangas = utils.GetSlicedResults[models.Manga](mangasChan)

	return mangas, nil
}
