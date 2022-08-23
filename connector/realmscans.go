package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type realm models.Source

func (r *realm) GetSource() models.Source {
	return models.Source(*r)
}

func (r *realm) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	sels := models.MangaListSelectors{
		URL:                   "https://realmscans.com/series",
		MangaTitleSelector:    "div.listupd > div.bs div.tt",
		MangaImageURLSelector: "div.listupd > div.bs img[data-src,src]",
		MangaURLSelector:      "div.listupd > div.bs a[href]",
		NextPageSelector:      "div.hpage a.r[href]",
	}

	return scrapper.ScrapeMangaList(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: r.Transport})
}

func (r *realm) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	sels := models.MangaInfoSelectors{
		URL:                       mangaURL,
		TitleSelector:             "div.info-right h1",
		ImageURLSelector:          "div.thumb img[src]",
		SynopsisSelector:          "div.info-right div.wd-full div.entry-content p",
		ChapterNumberSelector:     "div#chapterlist ul li[data-num]",
		ChapterTitleSelector:      "div#chapterlist ul li a span.chapternum",
		ChapterURLSelector:        "div#chapterlist ul li a[href]",
		ChapterUploadDateSelector: "div#chapterlist ul li a span.chapterdate",
		ChapterUploadDateFormat:   "January 2, 2006",
	}

	return scrapper.ScrapeMangaInfo(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: r.Transport})
}

func (r *realm) GetChapterPages(ctx web.Context, chapterInfoUrl string) ([]string, error) {
	sels := models.ChapterInfoSelectors{
		URL:          chapterInfoUrl,
		PageSelector: "#readerarea img[data-src,src]",
	}

	data, err := scrapper.ScrapeChapterPages(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: r.Transport})
	if err != nil || len(data) == 0 {
		injScript := scrapper.GetInjectionScript(sels.PageSelector)
		return scrapper.SimulateBrowser(ctx, chapterInfoUrl, injScript)
	}

	return data, err
}

func getRealmScansConnector() models.IConnector {
	return &realm{
		Name:      "Realm Scans",
		Domain:    "realmscans.com",
		IconURL:   "https://cdn.realmscans.com/2021/09/logo-realm-scans-2.webp",
		Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
	}
}
