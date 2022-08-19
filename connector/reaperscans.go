package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type reaper models.Source

func (r *reaper) Get() models.Source {
	return models.Source(*r)
}

func (r *reaper) GetDomain() string {
	return r.Domain
}

func (r *reaper) GetName() string {
	return r.Name
}

func (r *reaper) GetIconURL() string {
	return r.IconURL
}

func (r *reaper) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	sels := models.MangaListSelectors{
		URL:                   "https://reaperscans.com/all-series/comics/",
		MangaTitleSelector:    "div.manga h3 a",
		MangaImageURLSelector: "div.listupd > div.bs img[src]",
		MangaURLSelector:      "div.manga h3 a[href]",
		NextPageSelector:      "div.hpage r.r[href]",
	}

	return scrapper.ScrapeMangaList(ctx, sels, r.Transport)
}

func (r *reaper) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
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

	return scrapper.ScrapeMangaInfo(ctx, sels, r.Transport)
}

func (r *reaper) GetChapterPages(ctx web.Context, chapterInfoUrl string) ([]string, error) {
	sels := models.ChapterInfoSelectors{
		URL:          chapterInfoUrl,
		PageSelector: "div#readerarea p img[src]",
	}

	return scrapper.ScrapeChapterPages(ctx, sels, r.Transport)
}

func getreaperScansConnector() models.IConnector {
	return &reaper{
		Name:      "reaper Scans",
		Domain:    "reaperscans.com",
		IconURL:   "https://www.reaperscans.com/wp-content/uploads/2021/03/Group_1.png",
		Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
	}
}
