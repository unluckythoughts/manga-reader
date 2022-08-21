package scrapper

import (
	"io"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func populateMangas(ctx web.Context, sels models.MangaListSelectors, resp *mangaListResponse) func(h *colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		var titles []string
		titles, resp.Error = getTextListForSelector(h, sels.MangaTitleSelector)
		if resp.Error != nil {
			return
		}

		var urls []string
		urls, resp.Error = getTextListForSelector(h, sels.MangaURLSelector)
		if resp.Error != nil {
			return
		}

		var imageUrls []string
		imageUrls, resp.Error = getTextListForSelector(h, sels.MangaImageURLSelector)
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
		resp.NextPage, resp.Error = getText(h.DOM.Find(sels.NextPageSelector), sels.NextPageSelector)
	}
}

type mangaListResponse struct {
	Mangas   []models.Manga
	NextPage string
	Error    error
}

type ScrapeOptions struct {
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

func ScrapeMangaList(ctx web.Context, sels models.MangaListSelectors, opts *ScrapeOptions) ([]models.Manga, error) {
	opts.SetDefaults()

	resp := mangaListResponse{
		NextPage: sels.URL,
	}

	c := getColly(ctx, opts.RoundTripper)
	c.OnHTML(opts.InitialHtmlTag, populateMangas(ctx, sels, &resp))

	for resp.NextPage != "" {
		if isInternalLink(resp.NextPage) {
			resp.NextPage = sels.URL + resp.NextPage
		}

		err := c.Request(opts.RequestMethod, resp.NextPage, opts.Body, nil, opts.Headers)
		if err != nil {
			return resp.Mangas, err
		}

		if resp.Error != nil {
			return resp.Mangas, resp.Error
		}
	}

	return resp.Mangas, resp.Error
}
