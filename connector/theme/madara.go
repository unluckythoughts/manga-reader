package theme

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

func GetMadaraScrapeOptsForMangaList(c models.Connector) scrapper.ScrapeOptions {
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

	return opts
}

func GetMadaraScrapeOptsForChapterList(c models.Connector, manga_id, chaptersURL string) scrapper.ScrapeOptions {
	headers := http.Header{}
	headers.Set("content-type", "application/x-www-form-urlencoded")
	headers.Set("referer", c.BaseURL)

	params := url.Values{}
	params.Add("action", "manga_get_chapters")
	params.Add("manga", manga_id)

	opts := scrapper.ScrapeOptions{
		URL:            c.BaseURL + c.MangaListPath,
		RoundTripper:   c.Transport,
		RequestMethod:  http.MethodPost,
		InitialHtmlTag: scrapper.WHOLE_BODY_TAG,
		Headers:        headers,
		Body:           strings.NewReader(params.Encode()),
	}
	if chaptersURL != "" {
		opts.URL = chaptersURL
	}

	return opts
}
