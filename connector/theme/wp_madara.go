package theme

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

func GetMadaraScrapeOptsForMangaList(c models.Connector, page string) scrapper.ScrapeOptions {
	headers := http.Header{}
	headers.Set("content-type", "application/x-www-form-urlencoded")
	headers.Set("referer", c.BaseURL+c.MangaListPath)

	params := url.Values{}
	params.Add("action", "madara_load_more")
	params.Add("template", "madara-core/content/content-archive")
	params.Add("page", page)
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

type madara models.Connector

func GetMadaraConnector() *madara {
	return &madara{
		MangaListPath: "wp-admin/admin-ajax.php",
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		Selectors: models.Selectors{
			List: models.MangaList{
				MangaContainer: "div.page-item-detail.manga",
				MangaTitle:     "h3 a",
				MangaImageURL:  "img[src],img[data-src]",
				MangaURL:       "h3 a[href]",
				NextPage:       "",
			},
			Info: models.MangaInfo{
				Title:                   ".post-title h1",
				ImageURL:                ".profile-manga .summary_image a img[data-src],.profile-manga .summary_image a img[src]",
				OtherID:                 ".add-bookmark a[data-post]",
				Synopsis:                ".summary__content p:last-of-type",
				ChapterContainer:        "ul.main li",
				ChapterNumber:           "a",
				ChapterTitle:            "a",
				ChapterURL:              "a[href]",
				ChapterUploadDate:       "a+span i",
				ChapterUploadDateFormat: "January 2, 2006",
			},
			Chapter: models.PageSelectors{
				ImageUrl: ".reading-content img.wp-manga-chapter-img[data-src],.reading-content img.wp-manga-chapter-img[src]",
			},
		},
	}
}

func (m *madara) GetSource() models.Source {
	return m.Source
}

func (m *madara) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*m)

	var mangas []models.Manga
	for i := 0; true; i++ {
		opts := GetMadaraScrapeOptsForMangaList(c, strconv.Itoa(i))
		pageMangas, err := scrapper.ScrapeMangas(ctx, c, &opts)
		if err != nil {
			return mangas, err
		}

		mangas = append(mangas, pageMangas...)

		if len(pageMangas) <= 0 {
			break
		}
	}

	return mangas, nil
}

func (m *madara) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*m)
	opts := scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	manga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
	if err != nil {
		return manga, err
	}

	if len(manga.Chapters) == 0 {
		opts = GetMadaraScrapeOptsForChapterList(c, manga.OtherID, mangaURL+"ajax/chapters")
		chaptersManga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
		if err != nil {
			return manga, err
		}

		manga.Chapters = chaptersManga.Chapters
	}

	return manga, err
}

func (m *madara) GetChapterPages(ctx web.Context, chapterUrl string) (models.Pages, error) {
	c := models.Connector(*m)
	opts := scrapper.ScrapeOptions{
		URL:          chapterUrl,
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeChapterPages(ctx, c, &opts)
}
