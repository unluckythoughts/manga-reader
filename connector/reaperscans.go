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
	listURL := "https://reaperscans.com/wp-admin/admin-ajax.php"
	sels := models.MangaListSelectors{
		URL:                   listURL,
		MangaTitleSelector:    "div.page-item-detail.manga h3 a",
		MangaImageURLSelector: "div.page-item-detail.manga img[data-src,src]",
		MangaURLSelector:      "div.page-item-detail.manga h3 a[href]",
		NextPageSelector:      "",
	}

	headers := http.Header{}
	headers.Set("content-type", "application/x-www-form-urlencoded")
	headers.Set("referer", listURL)

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
		RoundTripper:   r.Transport,
		RequestMethod:  http.MethodPost,
		Headers:        headers,
		InitialHtmlTag: scrapper.WHOLE_BODY_TAG,
		Body:           strings.NewReader(params.Encode()),
	}

	return scrapper.ScrapeMangaList(ctx, sels, &opts)
}

func (r *reaper) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	sels := models.MangaInfoSelectors{
		URL:                       mangaURL,
		TitleSelector:             "div.container .post-title h1",
		ImageURLSelector:          "div.tab-summary a img[data-src,src], a#roi img[data-src,src], a#roiroi img[data-src,src]",
		SynopsisSelector:          "div.container .summary__content",
		ChapterNumberSelector:     "div.listing-chapters_wrap ul.main li .chapter-link a > p",
		ChapterTitleSelector:      "div.listing-chapters_wrap ul.main li .chapter-link a > p",
		ChapterURLSelector:        "div.listing-chapters_wrap ul.main li .chapter-link a[href]",
		ChapterUploadDateSelector: "div.listing-chapters_wrap ul.main li .chapter-link a span i",
		ChapterUploadDateFormat:   "Jan 02, 2006",
	}

	return scrapper.ScrapeMangaInfo(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: r.Transport})
}

func (r *reaper) GetChapterPages(ctx web.Context, chapterInfoUrl string) ([]string, error) {
	sels := models.ChapterInfoSelectors{
		URL:          chapterInfoUrl,
		PageSelector: "div.reading-content img.wp-manga-chapter-img[data-src,src]",
	}

	return scrapper.ScrapeChapterPages(ctx, sels, &scrapper.ScrapeOptions{RoundTripper: r.Transport})
}

func getreaperScansConnector() models.IConnector {
	return &reaper{
		Name:      "reaper Scans",
		Domain:    "reaperscans.com",
		IconURL:   "https://reaperscans.com/wp-content/uploads/2021/07/logo-reaper-2.png",
		Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
	}
}
