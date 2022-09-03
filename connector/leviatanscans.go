package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type leviatan models.Connector

func GetLeviatanScansConnector() models.IConnector {
	return &leviatan{
		Source: models.Source{
			Name:    "Leviatan Scans",
			Domain:  "leviatanscans.com",
			IconURL: "https://styles.redditmedia.com/t5_2hfywp/styles/communityIcon_qdo3swk6vzl41.png",
		},
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		BaseURL:       "https://leviatanscans.com/",
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
				Title:                   "#manga-title h1",
				ImageURL:                ".profile-manga .summary_image a img[data-src],.profile-manga .summary_image a img[src]",
				OtherID:                 "#manga-chapters-holder[data-id]",
				Synopsis:                ".summary_content .post-content_item:last-of-type p span",
				ChapterContainer:        "ul.main li",
				ChapterNumber:           "a",
				ChapterTitle:            "a",
				ChapterURL:              "a[href]",
				ChapterUploadDate:       "a+span i",
				ChapterUploadDateFormat: "Jan 2, 2006",
			},
			Chapter: models.PageSelectors{
				ImageUrl: ".reading-content img.wp-manga-chapter-img[data-src],.reading-content img.wp-manga-chapter-img[src]",
			},
		},
	}
}

func (r *leviatan) GetSource() models.Source {
	return r.Source
}

func (r *leviatan) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*r)
	opts := theme.GetMadaraScrapeOptsForMangaList(c)

	return scrapper.ScrapeMangas(ctx, c, &opts)
}

func (r *leviatan) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*r)
	opts := scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	manga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
	if err != nil {
		return manga, err
	}

	if len(manga.Chapters) == 0 {
		opts = theme.GetMadaraScrapeOptsForChapterList(c, manga.OtherID, mangaURL+"ajax/chapters")
		chaptersManga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
		if err != nil {
			return manga, err
		}

		manga.Chapters = chaptersManga.Chapters
	}

	return manga, err
}

func (r *leviatan) GetChapterPages(ctx web.Context, chapterUrl string) (models.Pages, error) {
	c := models.Connector(*r)
	opts := scrapper.ScrapeOptions{
		URL:          chapterUrl,
		RoundTripper: c.Transport,
	}
	opts.SetDefaults()
	return scrapper.ScrapeChapterPages(ctx, c, &opts)
}
