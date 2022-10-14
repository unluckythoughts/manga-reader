package connector

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type leviatan struct {
	c    models.IMangaConnector
	conn models.MangaConnector
}

func GetLeviatanScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Leviatan Scans",
		Domain:  "leviatanscans.com",
		IconURL: "https://en.leviatanscans.com/wp-content/uploads/2022/08/cropped-isotiponegro-32x32.png",
	}

	c.BaseURL = "https://en.leviatanscans.com/"

	c.MangaSelectors.List.MangaContainer = ".page-content-listing .manga"
	c.MangaSelectors.List.MangaTitle = ".item-summary h3 a,.item-summary h5 a"
	c.MangaSelectors.List.MangaURL = ".item-summary h3 a[href],.item-summary h5 a[href]"
	c.MangaSelectors.List.MangaImageURL = ".item-thumb a img[data-src],.item-thumb a img[src]"
	c.MangaSelectors.List.LastPage = ".wp-pagenavi a.last"
	c.MangaSelectors.List.PageParam = "/page/" + scrapper.PAGE_ID

	c.MangaSelectors.Info.Title = "#manga-title h1"
	c.MangaSelectors.Info.ImageURL = ".profile-manga .summary_image img[data-src], .profile-manga .summary_image img[src]"
	c.MangaSelectors.Info.Synopsis = ".profile-manga .post-content_item:last-of-type p"
	c.MangaSelectors.Info.ChapterContainer = "ul.main li"
	c.MangaSelectors.Info.ChapterNumber = "a"
	c.MangaSelectors.Info.ChapterTitle = "a"
	c.MangaSelectors.Info.ChapterURL = "a[href]"
	c.MangaSelectors.Info.ChapterUploadDate = "span.chapter-release-date"
	c.MangaSelectors.Info.ChapterUploadDateFormat = "Jan 2, 2006"

	c.MangaSelectors.Chapter.ImageUrl = ".reading-content img[data-src],.reading-content img[src]"

	return &leviatan{c: c, conn: models.MangaConnector(*c)}
}

func (l *leviatan) GetSource() models.Source {
	return l.c.GetSource()
}

func (l *leviatan) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	return l.c.GetMangaList(ctx)
}

func (l *leviatan) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
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

func (l *leviatan) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	return l.c.GetChapterPages(ctx, chapterURL)
}
