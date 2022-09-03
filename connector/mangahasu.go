package connector

import (
	"net/http"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type mangahasu models.Connector

func GetMangaHasuConnector() models.IConnector {
	return &mangahasu{
		Source: models.Source{
			Name:    "Manga Hasu",
			Domain:  "mangahasu.se",
			IconURL: "https://mangahasu.se/favicon.ico",
		},
		BaseURL:       "http://mangahasu.se/",
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		MangaListPath: "directory.html",
		Selectors: models.Selectors{
			List: models.MangaList{
				MangaContainer: ".list_manga .div_item",
				MangaTitle:     ".info-manga a h3",
				MangaURL:       ".info-manga a[href]",
				MangaImageURL:  ".wrapper_imagage img[src],img[src],.wrapper_imagage a[src]",
				NextPage:       ".pagination-ct a[title='Tiếp']",
				LastPage:       ".pagination-ct a[title='Trang cuối']",
			},
			Info: models.MangaInfo{
				Title:                   ".wrapper_content .info-title h1",
				ImageURL:                ".wrapper_content .info-img img[src]",
				Synopsis:                ".wrapper_content .content-info > h3 + div",
				ChapterContainer:        ".wrapper_content .list-chapter tbody tr",
				ChapterNumber:           "td.name a",
				ChapterTitle:            "td.name a",
				ChapterURL:              "td.name a[href]",
				ChapterUploadDate:       "td.date-updated",
				ChapterUploadDateFormat: "Jan 02, 2006",
			},
			Chapter: models.PageSelectors{
				ImageUrl: "#loadchapter .img img[data-src], #loadchapter .img img[src]",
			},
		},
	}
}

func (m *mangahasu) GetSource() models.Source {
	return m.Source
}

func (m *mangahasu) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*m)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.MangaListPath,
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeMangasParallel(ctx, c, opts)
}

func (m *mangahasu) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*m)
	opts := &scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	manga, err := scrapper.ScrapeMangaInfo(ctx, c, opts)

	for i, c := range manga.Chapters {
		texts := strings.Split(c.Title, "</span>")
		if len(texts) > 1 {
			manga.Chapters[i].Title = scrapper.GetChapterTitle(texts[1])
			manga.Chapters[i].Number = scrapper.GetChapterNumber(texts[1])
		}
	}

	return manga, err
}

func (m *mangahasu) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	c := models.Connector(*m)

	headers := http.Header{}
	headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	opts := &scrapper.ScrapeOptions{
		URL:          chapterURL,
		RoundTripper: c.Transport,
		Headers:      headers,
	}
	return scrapper.ScrapeChapterPages(ctx, c, opts)
}