package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type realm models.Connector

func GetRealmScansConnector() models.IConnector {
	return &realm{
		Source: models.Source{
			Name:      "Realm Scans",
			Domain:    "realmscans.com",
			IconURL:   "https://cdn.realmscans.com/2021/09/logo-realm-scans-2.webp",
			Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		},
		BaseURL:       "https://realmscans.com/",
		MangaListPath: "series",
		Selectors: models.Selectors{
			List: models.MangaList{
				MangaContainer: "div.listupd > div.bs",
				MangaTitle:     ".tt",
				MangaURL:       "img[data-src], img[src]",
				MangaImageURL:  "a[href]",
				NextPage:       "div.hpage a.r[href]",
			},
			Info: models.MangaInfo{
				Title:                   ".info-right h1",
				ImageURL:                ".thumb img[src]",
				Synopsis:                ".info-right .wd-full .entry-content p",
				ChapterContainer:        "div#chapterlist ul li",
				ChapterNumber:           "[data-num]",
				ChapterTitle:            "a span.chapternum",
				ChapterURL:              "a[href]",
				ChapterUploadDate:       "a span.chapterdate",
				ChapterUploadDateFormat: "January 2, 2006",
			},
			Chapter: models.PageSelectors{
				ImageUrl: "#readerarea img[data-src,src]",
			},
		},
	}
}

func (r *realm) GetSource() models.Source {
	return r.Source
}

func (r *realm) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*r)
	opts := &scrapper.ScrapeOptions{
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.ScrapeMangas(ctx, c, opts)
}

func (r *realm) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*r)
	opts := &scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.NewScrapeMangaInfo(ctx, c, opts)
}

func (r *realm) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	c := models.Connector(*r)
	opts := &scrapper.ScrapeOptions{
		URL:          chapterURL,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	pages, err := scrapper.NewScrapeChapterPages(ctx, c, opts)
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
