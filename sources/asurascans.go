package sources

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetAsuraScansSource() models.Source {
	return models.Source{
		Name:         "Asura Scans",
		IconURL:      "https://www.asurascans.com/wp-content/uploads/2021/03/Group_1.png",
		RoundTripper: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		MangaList: models.MangaListSelectors{
			URL:                   "http://asurascans.com/manga",
			MangaTitleSelector:    "div.listupd > div.bs div.tt",
			MangaImageURLSelector: "div.listupd > div.bs img[src]",
			MangaURLSelector:      "div.listupd > div.bs a[href]",
			NextPageSelector:      "div.hpage a.r[href]",
		},
		MangaInfo: models.MangaInfoSelectors{
			TitleSelector:             "div.infox > h1",
			ImageURLSelector:          "div.thumb img[src]",
			SynopsisSelector:          "div.infox > div.wd-full div.entry-content p",
			ChapterNumberSelector:     "ul.clstyle li[data-num]",
			ChapterTitleSelector:      "ul.clstyle li a span.chapternum",
			ChapterURLSelector:        "ul.clstyle li a[href]",
			ChapterUploadDateSelector: "ul.clstyle li a span.chapterdate",
		},
		ChapterInfo: models.ChapterInfoSelectors{
			ImageURLsSelector: "div#readerarea p img[src]",
		},
	}
}
