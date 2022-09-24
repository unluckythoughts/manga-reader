package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type lnp models.NovelConnector

func GetLightNovelPubConnector() models.INovelConnector {
	c := &lnp{
		Source: models.Source{
			Name:    "Light Novel Pub",
			Domain:  "lightnovelpub.com",
			IconURL: "https://static.lightnovelpub.com/content/img/lightnovelpub/logo.png",
		},
		BaseURL:       "https://www.lightnovelpub.com/",
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		NovelListPath: "stories-17091737/genre-all/order-new/status-all/p-1",
		NovelSelectors: models.NovelSelectors{
			List: models.NovelList{
				NovelContainer: "ul.novel-list li",
				NovelTitle:     ".item-body h4.novel-title a",
				NovelURL:       ".item-body h4.novel-title a[href]",
				NovelImageURL:  ".cover-wrap img[data-src], .cover-wrap img[src]",
				NextPage:       "ul.pagination li.PagedList-skipToNext a[href]",
			},
			Info: models.NovelInfo{
				Title:                   ".novel-info h1.novel-title",
				ImageURL:                "figure.cover img[data-src], figure.cover img[src]",
				Synopsis:                ".novel-body .summary .content p",
				ChapterListURL:          "chapters",
				ChapterListNextPage:     "ul.pagination li.PagedList-skipToNext a[href]",
				ChapterContainer:        "ul.chapter-list li",
				ChapterNumber:           "a span.chapter-no",
				ChapterTitle:            "a strong.chapter-title",
				ChapterUploadDate:       "a time.chapter-update[datetime]",
				ChapterUploadDateFormat: "2006-01-02 03:04",
				ChapterURL:              "a[href]",
			},
			Chapter: models.NovelChapterTextSelectors{
				Paragraph: "#chapter-container p",
			},
		},
	}

	return c
}

func (n *lnp) GetSource() models.Source {
	return n.Source
}

func (n *lnp) GetNovelList(ctx web.Context) ([]models.Novel, error) {
	c := models.NovelConnector(*n)
	opts := &scrapper.ScrapeOptions{
		URL:          c.BaseURL + c.NovelListPath,
		RoundTripper: c.Transport,
	}

	return scrapper.ScrapeNovels(ctx, c, opts)
}

func (n *lnp) GetNovelInfo(ctx web.Context, novelURL string) (models.Novel, error) {
	c := models.NovelConnector(*n)
	opts := &scrapper.ScrapeOptions{
		URL:          novelURL,
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeNovelInfo(ctx, c, opts)
}

func (n *lnp) GetNovelChapter(ctx web.Context, chapterURL string) ([]string, error) {
	c := models.NovelConnector(*n)

	opts := &scrapper.ScrapeOptions{
		URL:          chapterURL,
		RoundTripper: c.Transport,
	}

	return scrapper.ScrapeNovelChapterText(ctx, c, opts)
}
