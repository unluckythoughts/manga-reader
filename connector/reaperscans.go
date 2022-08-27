package connector

import (
	"net/http"
	"net/url"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type reaper models.Connector

func GetReaperScansConnector() models.IConnector {
	return &reaper{
		Source: models.Source{
			Name:    "Reaper Scans",
			Domain:  "reaperscans.com",
			IconURL: "https://styles.redditmedia.com/t5_4zgiee/styles/communityIcon_gxpzm2tt41l71.png",
		},
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		BaseURL:       "https://reaperscans.com/",
		MangaListPath: "wp-admin/admin-ajax.php",
		Selectors: models.Selectors{
			List: models.MangaList{
				MangaContainer: "div.page-item-detail.manga",
				MangaTitle:     "h3 a",
				MangaImageURL:  "img[data-src],img[src]",
				MangaURL:       "h3 a[href]",
				NextPage:       "",
			},
			Info: models.MangaInfo{
				Title:                   ".container .post-title h1",
				ImageURL:                ".tab-summary a img[data-src],.tab-summary a img[src],a#roi img[data-src],a#roi img[src],a#roiroi img[data-src],a#roiroi img[src]",
				Synopsis:                ".container .summary__content",
				ChapterContainer:        ".listing-chapters_wrap ul.main li",
				ChapterNumber:           ".chapter-link a > p",
				ChapterTitle:            ".chapter-link a > p",
				ChapterURL:              ".chapter-link a[href]",
				ChapterUploadDate:       ".chapter-link a span i",
				ChapterUploadDateFormat: "Jan 02, 2006",
			},
			Chapter: models.PageSelectors{
				ImageUrl: ".reading-content img.wp-manga-chapter-img[data-src],.reading-content img.wp-manga-chapter-img[src]",
			},
		},
	}
}

func (r *reaper) GetSource() models.Source {
	return r.Source
}

func (r *reaper) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*r)
	headers := http.Header{}
	headers.Set("content-type", "application/x-www-form-urlencoded")
	headers.Set("referer", c.BaseURL+c.MangaListPath)

	params := url.Values{}
	params.Add("action", "madara_load_more")
	params.Add("template", "madara-core/content/content-archive")
	params.Add("page", "0")
	params.Add("vars[paged]", "1")
	params.Add("vars[orderby]", "post_title")
	params.Add("vars[template]", "archive")
	params.Add("vars[sidebar]", "full")
	params.Add("vars[meta_query][0][0][key]", "_wp_manga_chapter_type")
	params.Add("vars[meta_query][0][0][value]", "manga")
	params.Add("vars[meta_query][0][relation]", "AND")
	params.Add("vars[meta_query][relation]", "OR")
	params.Add("vars[post_type]", "wp-manga")
	params.Add("vars[order]", "ASC")
	params.Add("vars[posts_per_page]", "500")

	opts := scrapper.ScrapeOptions{
		URL:            c.BaseURL + c.MangaListPath,
		RoundTripper:   c.Transport,
		RequestMethod:  http.MethodPost,
		InitialHtmlTag: scrapper.WHOLE_BODY_TAG,
		Headers:        headers,
		Body:           strings.NewReader(params.Encode()),
	}
	opts.SetDefaults()
	return scrapper.ScrapeMangas(ctx, c, &opts)
}

func (r *reaper) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*r)
	opts := scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.ScrapeMangaInfo(ctx, c, &opts)
}

func (r *reaper) GetChapterPages(ctx web.Context, chapterUrl string) (models.Pages, error) {
	c := models.Connector(*r)
	opts := scrapper.ScrapeOptions{
		URL:          chapterUrl,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.ScrapeChapterPages(ctx, c, &opts)
}
