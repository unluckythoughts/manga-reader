package connector

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type astra struct {
	c    models.IMangaConnector
	conn models.MangaConnector
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

	c.MangaSelectors.List.MangaContainer = ".manga"
	c.MangaSelectors.List.MangaTitle = ".post-title h3 a,.post-title .h5 a"
	c.MangaSelectors.List.MangaURL = ".post-title h3 a[href],.post-title .h5 a[href]"
	c.MangaSelectors.List.MangaImageURL = ".item-thumb a img[data-src],.item-thumb a img[src]"

	c.MangaSelectors.Info.Title = ".manga-excerpt h1"
	c.MangaSelectors.Info.ImageURL = ".profile-manga .summary_image img[data-src], .profile-manga .summary_image img[src]"
	c.MangaSelectors.Info.Synopsis = ".manga-excerpt"
	c.MangaSelectors.Info.ChapterContainer = "ul.main li"
	c.MangaSelectors.Info.ChapterNumber = "a"
	c.MangaSelectors.Info.ChapterTitle = "a"
	c.MangaSelectors.Info.ChapterURL = "a[href]"
	c.MangaSelectors.Info.ChapterUploadDate = "span.chapter-release-date,span.chapter-release-date a[title]"
	c.MangaSelectors.Info.ChapterUploadDateFormat = "January 2, 2006"

	c.MangaSelectors.Chapter.ImageUrl = ".reading-content img[data-src],.reading-content img[src]"

	return &astra{c: c, conn: models.MangaConnector(*c)}
}

func (l *astra) GetSource() models.Source {
	return l.c.GetSource()
}

func (l *astra) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	return l.c.GetMangaList(ctx)
}

func (l *astra) GetLatestMangaList(ctx web.Context, latestTitle string) ([]models.Manga, error) {
	return l.c.GetLatestMangaList(ctx, latestTitle)
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
