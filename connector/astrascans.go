package connector

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type astra struct {
	c    models.IMangaConnector
	conn models.Connector
}

func GetAstraScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Astra Scans",
		Domain:  "astrascans.com",
		IconURL: "https://astrascans.com/wp-content/uploads/2022/06/cropped-logo.png",
	}

	c.BaseURL = "https://astrascans.com/"
	c.MangaListPath = "manga/"

	c.Selectors.List.MangaContainer = ".manga"
	c.Selectors.List.MangaTitle = ".post-title h3 a,.post-title .h5 a"
	c.Selectors.List.MangaURL = ".post-title h3 a[href],.post-title .h5 a[href]"
	c.Selectors.List.MangaImageURL = ".item-thumb a img[data-src],.item-thumb a img[src]"

	c.Selectors.Info.Title = ".manga-excerpt h1"
	c.Selectors.Info.ImageURL = ".profile-manga .summary_image img[data-src], .profile-manga .summary_image img[src]"
	c.Selectors.Info.Synopsis = ".manga-excerpt"
	c.Selectors.Info.ChapterContainer = "ul.main li"
	c.Selectors.Info.ChapterNumber = "a"
	c.Selectors.Info.ChapterTitle = "a"
	c.Selectors.Info.ChapterURL = "a[href]"
	c.Selectors.Info.ChapterUploadDate = "span.chapter-release-date,span.chapter-release-date a[title]"
	c.Selectors.Info.ChapterUploadDateFormat = "January 2, 2006"

	c.Selectors.Chapter.ImageUrl = ".reading-content img[data-src],.reading-content img[src]"

	return &astra{c: c, conn: models.Connector(*c)}
}

func (l *astra) GetSource() models.Source {
	return l.c.GetSource()
}

func (l *astra) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	return l.c.GetMangaList(ctx)
}

func (l *astra) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	manga, err := l.c.GetMangaInfo(ctx, mangaURL)
	if err != nil {
		return manga, err
	}

	if len(manga.Chapters) == 0 {
		c := l.conn
		opts := theme.GetMadaraScrapeOptsForChapterList(c, manga.OtherID, mangaURL+"ajax/chapters/")
		chaptersManga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
		if err != nil {
			return manga, err
		}

		manga.Chapters = chaptersManga.Chapters
	}

	return manga, err
}

func (l *astra) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	return l.c.GetChapterPages(ctx, chapterURL)
}
