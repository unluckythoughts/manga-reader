package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type flame models.Connector

func GetFlameScansConnector() models.IConnector {
	return &flame{
		Source: models.Source{
			Name:    "Flame Scans",
			Domain:  "flamescans.org",
			IconURL: "https://flamescans.org/favicon.ico",
		},
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		BaseURL:       "http://flamescans.org/",
		MangaListPath: "series/",
		Selectors: models.Selectors{
			List: models.MangaList{
				MangaContainer: ".listupd > .bs",
				MangaTitle:     ".tt",
				MangaImageURL:  "img[src]",
				MangaURL:       "a[href]",
				NextPage:       ".hpage a.r[href]",
			},
			Info: models.MangaInfo{
				Title:                   ".info-half h1.entry-title",
				ImageURL:                ".thumb-half .thumb img[src]",
				Synopsis:                ".info-half .summary .wd-full .entry-content",
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

func (f *flame) GetSource() models.Source {
	return f.Source
}

func (f *flame) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*f)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.MangaListPath,
		RoundTripper: c.Transport,
	}
		return scrapper.ScrapeMangas(ctx, c, opts)
}

func (f *flame) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*f)
	opts := &scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
		return scrapper.ScrapeMangaInfo(ctx, c, opts)
}

func (f *flame) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	c := models.Connector(*f)
	opts := &scrapper.ScrapeOptions{
		URL:          chapterURL,
		RoundTripper: c.Transport,
	}
		return scrapper.ScrapeChapterPages(ctx, c, opts)
}
