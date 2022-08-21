package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type asura models.Source

func (a *asura) Get() models.Source {
	return models.Source(*a)
}

func (a *asura) GetDomain() string {
	return a.Domain
}

func (a *asura) GetName() string {
	return a.Name
}

func (a *asura) GetIconURL() string {
	return a.IconURL
}

func (a *asura) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	sels := models.MangaListSelectors{
		URL:                   "http://asurascans.com/manga",
		MangaTitleSelector:    "div.listupd > div.bs div.tt",
		MangaImageURLSelector: "div.listupd > div.bs img[src]",
		MangaURLSelector:      "div.listupd > div.bs a[href]",
		NextPageSelector:      "div.hpage a.r[href]",
	}

	return scrapper.ScrapeMangaList(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: a.Transport})
}

func (a *asura) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	sels := models.MangaInfoSelectors{
		URL:                       mangaURL,
		TitleSelector:             "div.infox > h1",
		ImageURLSelector:          "div.thumb img[src]",
		SynopsisSelector:          "div.infox > div.wd-full div.entry-content p",
		ChapterNumberSelector:     "ul.clstyle li[data-num]",
		ChapterTitleSelector:      "ul.clstyle li a span.chapternum",
		ChapterURLSelector:        "ul.clstyle li a[href]",
		ChapterUploadDateSelector: "ul.clstyle li a span.chapterdate",
	}

	return scrapper.ScrapeMangaInfo(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: a.Transport})
}

func (a *asura) GetChapterPages(ctx web.Context, chapterInfoUrl string) ([]string, error) {
	sels := models.ChapterInfoSelectors{
		URL:          chapterInfoUrl,
		PageSelector: "div#readerarea p img[src]",
	}

	return scrapper.ScrapeChapterPages(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: a.Transport})
}

func getAsuraScansConnector() models.IConnector {
	return &asura{
		Name:      "Asura Scans",
		Domain:    "asurascans.com",
		IconURL:   "https://www.asurascans.com/wp-content/uploads/2021/03/Group_1.png",
		Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
	}
}
