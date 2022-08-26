package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type asura models.Connector

func GetAsuraScansConnector() models.IConnector {
	return &asura{
		Source: models.Source{
			Name:      "Asura Scans",
			Domain:    "asurascans.com",
			IconURL:   "https://www.asurascans.com/wp-content/uploads/2021/03/Group_1.png",
			Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		},
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		BaseURL:       "http://asurascans.com/",
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
				Title:                   ".infox > h1",
				ImageURL:                ".thumb img[src]",
				Synopsis:                ".infox > .wd-full .entry-content p",
				ChapterContainer:        "#chapterlist ul.clstyle li",
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

func (a *asura) GetSource() models.Source {
	return a.Source
}

func (a *asura) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*a)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.MangaListPath,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.ScrapeMangas(ctx, c, opts)
}

func (a *asura) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*a)
	opts := &scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.NewScrapeMangaInfo(ctx, c, opts)
}

func (a *asura) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	c := models.Connector(*a)
	opts := &scrapper.ScrapeOptions{
		URL:          chapterURL,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.NewScrapeChapterPages(ctx, c, opts)
}
