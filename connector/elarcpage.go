package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type elarcpage models.Connector

func GetElarcPageConnector() models.IConnector {
	return &elarcpage{
		Source: models.Source{
			Name:    "Elarc Page",
			Domain:  "elarcpage.com",
			IconURL: "https://elarcpage.com/favicon.ico",
		},
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		BaseURL:       "http://elarcpage.com/",
		MangaListPath: "manga/",
		Selectors: models.Selectors{
			List: models.MangaList{
				MangaContainer: ".listupd > .bs",
				MangaTitle:     ".tt",
				MangaImageURL:  "img[src]",
				MangaURL:       "a[href]",
				NextPage:       ".hpage a.r[href]",
			},
			Info: models.MangaInfo{
				Title:                   ".info-right h1.entry-title",
				ImageURL:                ".info-left .thumb img[src]",
				Synopsis:                ".info-right .wd-full .entry-content p",
				ChapterContainer:        "#chapterlist ul li",
				ChapterNumber:           "[data-num]",
				ChapterTitle:            "a span.chapternum",
				ChapterURL:              "a[href]",
				ChapterUploadDate:       "a span.chapterdate",
				ChapterUploadDateFormat: "January 2, 2006",
			},
			Chapter: models.PageSelectors{
				ImageUrl: "#readerarea img[src]",
			},
		},
	}
}

func (e *elarcpage) GetSource() models.Source {
	return e.Source
}

func (e *elarcpage) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*e)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.MangaListPath,
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeMangas(ctx, c, opts)
}

func (e *elarcpage) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*e)
	opts := &scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeMangaInfo(ctx, c, opts)
}

func (e *elarcpage) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	c := models.Connector(*e)
	opts := &scrapper.ScrapeOptions{
		URL:          chapterURL,
		RoundTripper: c.Transport,
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
