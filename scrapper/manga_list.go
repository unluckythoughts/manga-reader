package scrapper

import (
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

const (
	WHOLE_BODY_TAG = "WHOLE_BODY_TAG"
)

type ScrapeOptions struct {
	URL            string
	RoundTripper   http.RoundTripper
	Headers        http.Header
	InitialHtmlTag string
	RequestMethod  string
	Body           io.Reader
}

func (opts *ScrapeOptions) SetDefaults() {
	if opts.InitialHtmlTag == "" {
		opts.InitialHtmlTag = "body"
	}

	if opts.RequestMethod == "" {
		opts.RequestMethod = http.MethodGet
	}
}

func ScrapeMangas(ctx web.Context, c models.Connector, opts *ScrapeOptions) ([]models.Manga, error) {
	mangas := []models.Manga{}
	// url := c.BaseURL + c.MangaListPath
	for opts.URL != "" {
		err := GetPageForScrapping(ctx, opts, func(h *colly.HTMLElement) {
			h.DOM.Find(c.Selectors.List.MangaContainer).Each(func(i int, s *goquery.Selection) {
				manga, err := GetMangaFromListSelectors(s, c.Selectors.List)
				if err != nil {
					ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
					return
				}

				if manga.Title != "" {
					mangas = append(mangas, manga)
				}

			})

			opts.URL = ""
			if nextPageElement, ok := GetElementForSelector(h.DOM, c.Selectors.List.NextPage); ok {
				if link, ok := nextPageElement.Attr("href"); ok {
					if strings.HasPrefix(link, c.MangaListPath) {
						opts.URL = c.BaseURL + link
					} else {
						opts.URL = c.BaseURL + c.MangaListPath + link
					}
				}
			}
		})

		if err != nil {
			ctx.Logger().With(zap.Error(err)).Debugf("error getting manga info from %s", c.Source.Domain)
			break
		}
	}

	return mangas, nil
}
