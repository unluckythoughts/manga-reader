package connector

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type aquamanga models.Connector

func GetAquaMangaConnector() models.IConnector {
	return &aquamanga{
		Source: models.Source{
			Name:    "Aqua Manga",
			Domain:  "aquamanga.com",
			IconURL: "https://aquamanga.com/wp-content/uploads/2021/03/cropped-cropped-favicon-1.png",
		},
		Transport:     cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		BaseURL:       "https://aquamanga.com/",
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
				Title:                   ".post-title h1",
				ImageURL:                ".profile-manga .summary_image a img[src]",
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
				ImageUrl: ".reading-content img.wp-manga-chapter-img[src]",
			},
		},
	}
}

func (a *aquamanga) GetSource() models.Source {
	return a.Source
}

func (a *aquamanga) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	c := models.Connector(*a)
	opts := theme.GetMadaraScrapeOptsForMangaList(c)

	return scrapper.ScrapeMangas(ctx, c, &opts)
}

func (a *aquamanga) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	c := models.Connector(*a)
	opts := scrapper.ScrapeOptions{
		URL:          mangaURL,
		RoundTripper: c.Transport,
	}
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

func (a *aquamanga) GetChapterPages(ctx web.Context, chapterUrl string) (models.Pages, error) {
	c := models.Connector(*a)
	opts := scrapper.ScrapeOptions{
		URL:          chapterUrl,
		RoundTripper: c.Transport,
	}
	return scrapper.ScrapeChapterPages(ctx, c, &opts)
}
