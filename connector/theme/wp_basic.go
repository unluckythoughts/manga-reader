package theme

import (
	"net/http"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type basic models.Connector

func GetBasicWordPressConnector() *basic {
	return &basic{
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		MangaListPath: "manga/",
		Selectors: models.Selectors{
			List: models.MangaList{
				MangaContainer: ".listupd > .bs",
				MangaTitle:     ".tt",
				MangaImageURL:  "img[data-src], img[src]",
				MangaURL:       "a[href]",
				NextPage:       ".hpage a.r[href]",
			},
			Info: models.MangaInfo{
				Title:                   ".infox > h1",
				ImageURL:                ".thumb img[src]",
				Synopsis:                ".infox > .wd-full .entry-content",
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

func (b *basic) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*b)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.MangaListPath,
		RoundTripper: c.Transport,
	}

	if c.List.LastPage != "" && strings.Contains(c.List.PageParam, scrapper.MANGA_LIST_PAGE_ID) {
		return scrapper.ScrapeMangasParallel(ctx, c, opts)
	}

	return scrapper.ScrapeMangas(ctx, c, opts)
}

func (b *basic) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*b)
	opts := &scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeMangaInfo(ctx, c, opts)
}

func (b *basic) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	c := models.Connector(*b)

	headers := http.Header{}
	headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	opts := &scrapper.ScrapeOptions{
		URL:          chapterURL,
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
