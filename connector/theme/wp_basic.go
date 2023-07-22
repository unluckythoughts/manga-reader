package theme

import (
	"net/http"
	"net/url"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type basic models.MangaConnector

func GetBasicWordPressConnector() *basic {
	return &basic{
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		MangaListPath: "manga/",
		MangaListURLParams: url.Values{
			"order": []string{"latest"},
		},
		MangaSelectors: models.MangaSelectors{
			List: models.MangaList{
				MangaContainer: ".listupd > .bs",
				MangaTitle:     ".tt",
				MangaImageURL:  "img[data-src], img[src]",
				MangaURL:       "a[href]",
				NextPage:       ".hpage a.r[href]",
			},
			Info: models.MangaInfo{
				Title:                   "h1.entry-title, h1",
				ImageURL:                ".thumb img[src]",
				Synopsis:                ".wd-full .entry-content",
				ChapterContainer:        "#chapterlist ul li",
				ChapterNumber:           "[data-num]",
				ChapterTitle:            "a span.chapternum",
				ChapterURL:              "a[href]",
				ChapterUploadDate:       "a span.chapterdate",
				ChapterUploadDateFormat: "January 2, 2006",
			},
			Chapter: models.PageSelectors{
				ImageUrl: "#readerarea p img[src]",
			},
		},
	}
}

func (b *basic) GetSource() models.Source {
	return b.Source
}

func (b *basic) GetMangaList(ctx web.Context) (mangas []models.Manga, err error) {
	c := models.MangaConnector(*b)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.MangaListPath,
		RoundTripper: c.Transport,
	}

	if c.List.LastPage != "" && strings.Contains(c.List.PageParam, scrapper.PAGE_ID) {
		mangas, err = scrapper.ScrapeMangasParallel(ctx, c, opts)
	} else {
		opts.URL = opts.URL + c.MangaListURLParams.Encode()
		mangas, err = scrapper.ScrapeMangas(ctx, c, opts)
	}

	for i := range mangas {
		mangas[i].URL = GetTrucattedURL(mangas[i].URL)
	}

	return mangas, err
}

func (b *basic) GetLatestMangaList(ctx web.Context, latestTitle string) (mangas []models.Manga, err error) {
	c := models.MangaConnector(*b)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.MangaListPath + c.MangaListURLParams.Encode(),
		RoundTripper: c.Transport,
	}

	mangas, err = scrapper.ScrapeLatestMangas(ctx, c, latestTitle, opts)

	for i := range mangas {
		mangas[i].URL = GetTrucattedURL(mangas[i].URL)
	}

	return mangas, err
}

func (b *basic) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.MangaConnector(*b)
	opts := &scrapper.ScrapeOptions{
		URL:          GetCompleteURL(mangaURL, b.Source.Domain),
		RoundTripper: c.Transport,
	}
	manga, err := scrapper.ScrapeMangaInfo(ctx, c, opts)
	manga.URL = GetTrucattedURL(manga.URL)

	for i := range manga.Chapters {
		manga.Chapters[i].URL = GetTrucattedURL(manga.Chapters[i].URL)
	}

	return manga, err
}

func (b *basic) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	c := models.MangaConnector(*b)

	headers := http.Header{}
	headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	opts := &scrapper.ScrapeOptions{
		URL:          GetCompleteURL(chapterURL, b.Source.Domain),
		RoundTripper: c.Transport,
		Headers:      headers,
	}

	pages, err := scrapper.ScrapeChapterPages(ctx, c, opts)
	if err != nil || len(pages.URLs) == 0 {
		injScript := scrapper.GetInjectionScript(c.Chapter.ImageUrl)
		imageURLs, err := scrapper.SimulateBrowser(ctx, chapterURL, injScript)
		if err != nil {
			return pages, err
		}

		pages.URLs = imageURLs
	}

	return pages, nil
}
