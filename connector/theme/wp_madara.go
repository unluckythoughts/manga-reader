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

func GetMadaraScrapeOptsForMangaList(c models.MangaConnector, page string) scrapper.ScrapeOptions {
	headers := http.Header{}
	headers.Set("content-type", "application/x-www-form-urlencoded")
	headers.Set("referer", c.BaseURL+c.MangaListPath)

	params := url.Values{}
	params.Add("action", "madara_load_more")
	params.Add("page", page)
	params.Add("template", "madara-core/content/content-archive")
	params.Add("vars[paged]", "1")
	params.Add("vars[orderby]", "date")
	params.Add("vars[sidebar]", "right")
	params.Add("vars[template]", "archive")
	params.Add("vars[post_type]", "wp-manga")
	params.Add("vars[post_status]", "publish")
	params.Add("vars[posts_per_page]", "500")
	params.Add("vars[manga_archives_item_layout]", "default")
	params.Add("vars[meta_query][0][paged]", "1")
	params.Add("vars[meta_query][0][orderby]", "date")
	params.Add("vars[meta_query][0][sidebar]", "right")
	params.Add("vars[meta_query][0][template]", "archive")
	params.Add("vars[meta_query][0][post_type]", "wp-manga")
	params.Add("vars[meta_query][0][post_status]", "publish")
	params.Add("vars[meta_query][relation]", "AND")

	// params.Add("vars[meta_query][0][0][key]", "_wp_manga_chapter_type")
	// params.Add("vars[meta_query][0][0][value]", "manga")
	// params.Add("vars[meta_query][0][relation]", "AND")

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

func GetMadaraScrapeOptsForChapterList(c models.MangaConnector, manga_id, chaptersURL string) scrapper.ScrapeOptions {
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

type madara models.MangaConnector

func GetMadaraConnector() *madara {
	return &madara{
		MangaListPath: "wp-admin/admin-ajax.php",
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		MangaSelectors: models.MangaSelectors{
			List: models.MangaList{
				MangaContainer: "div.page-item-detail.manga",
				MangaTitle:     "h3 a",
				MangaImageURL:  "img[data-src],img[src]",
				MangaURL:       "h3 a[href]",
				NextPage:       "",
			},
			Info: models.MangaInfo{
				Title:                   ".post-title h1",
				ImageURL:                ".profile-manga .summary_image a img[data-src],.profile-manga .summary_image a img[src]",
				OtherID:                 ".add-bookmark a[data-post]",
				Synopsis:                ".summary__content p:last-of-type",
				ChapterListURL:          "ajax/chapters",
				ChapterContainer:        "ul.main li",
				ChapterNumber:           "a",
				ChapterTitle:            "a",
				ChapterURL:              "a[href]",
				ChapterUploadDate:       "a+span i, a+span a[titlef]",
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
	c := models.MangaConnector(*m)

	var mangas []models.Manga
	for i := 0; true; i++ {
		opts := GetMadaraScrapeOptsForMangaList(c, strconv.Itoa(i))
		pageMangas, err := scrapper.ScrapeMangas(ctx, c, &opts)
		if err != nil {
			return mangas, err
		}

		for _, m := range pageMangas {
			m.URL = GetTrucattedURL(m.URL)
			mangas = append(mangas, m)
		}

		if len(pageMangas) <= 0 {
			break
		}
	}

	return mangas, nil
}

func (m *madara) GetLatestMangaList(ctx web.Context, latestTitle string) ([]models.Manga, error) {
	c := models.MangaConnector(*m)

	var mangas []models.Manga
	for i := 0; true; i++ {
		opts := GetMadaraScrapeOptsForMangaList(c, strconv.Itoa(i))
		pageMangas, err := scrapper.ScrapeMangas(ctx, c, &opts)
		if err != nil {
			return mangas, err
		}

		foundLast := false
		for _, m := range pageMangas {
			m.URL = GetTrucattedURL(m.URL)
			mangas = append(mangas, m)
			if m.Title == latestTitle {
				foundLast = true
				break
			}
		}

		if foundLast || len(pageMangas) <= 0 {
			break
		}
	}

	return mangas, nil
}

func (m *madara) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.MangaConnector(*m)
	opts := scrapper.ScrapeOptions{
		URL:          GetCompleteURL(mangaURL, m.Source.Domain),
		RoundTripper: c.Transport,
	}
	manga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
	if err != nil {
		return manga, err
	}

	manga.URL = GetTrucattedURL(manga.URL)
	if len(manga.Chapters) == 0 {
		chaptersURL := ""
		if c.Info.ChapterListURL != "" {
			chaptersURL = mangaURL + c.Info.ChapterListURL
		}

		opts = GetMadaraScrapeOptsForChapterList(c, manga.OtherID, chaptersURL)
		chaptersManga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
		if err != nil {
			return manga, err
		}

		manga.Chapters = chaptersManga.Chapters
	}

	for i := range manga.Chapters {
		manga.Chapters[i].URL = GetTrucattedURL(manga.Chapters[i].URL)
	}

	return manga, err
}

func (m *madara) GetChapterPages(ctx web.Context, chapterUrl string) (models.Pages, error) {
	c := models.MangaConnector(*m)
	opts := scrapper.ScrapeOptions{
		URL:          GetCompleteURL(chapterUrl, m.Source.Domain),
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeChapterPages(ctx, c, &opts)
}
