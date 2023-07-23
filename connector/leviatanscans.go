package connector

import (
	"net/url"

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
		IconURL: "/static/assets/images/leviatanscans.png",
	}

	c.BaseURL = "https://en.leviatanscans.com/"
	c.MangaListURLParams = url.Values{
		"m_orderby": []string{"latest"},
	}

	c.List.MangaContainer = ".page-content-listing .manga"
	c.List.MangaTitle = ".item-summary h3 a,.item-summary h5 a"
	c.List.MangaURL = ".item-summary h3 a[href],.item-summary h5 a[href]"
	c.List.MangaImageURL = ".item-thumb a img[data-src],.item-thumb a img[src]"
	c.List.LastPage = ".wp-pagenavi span.pages"
	c.List.PageParam = "/page/" + scrapper.PAGE_ID

	c.Info.Title = "#manga-title h1"
	c.Info.ImageURL = ".profile-manga .summary_image img[data-src], .profile-manga .summary_image img[src]"
	c.Info.Synopsis = ".profile-manga .post-content_item:last-of-type p"
	c.Info.ChapterContainer = "ul.main li"
	c.Info.ChapterNumber = "a"
	c.Info.ChapterTitle = "a"
	c.Info.ChapterURL = "a[href]"
	c.Info.ChapterUploadDate = "span.chapter-release-date"
	c.Info.ChapterUploadDateFormat = "Jan 2, 2006"

	c.Chapter.ImageUrl = ".reading-content img[data-src],.reading-content img[src]"

	return &leviatan{c: c, conn: models.MangaConnector(*c)}
}

func (l *leviatan) GetSource() models.Source {
	return l.c.GetSource()
}

func (l *leviatan) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	return l.c.GetMangaList(ctx)
}

func (l *leviatan) GetLatestMangaList(ctx web.Context, latestTitle string) ([]models.Manga, error) {
	return l.c.GetLatestMangaList(ctx, latestTitle)
}

func (l *leviatan) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	mangaURL = theme.GetCompleteURL(mangaURL, l.GetSource().Domain)
	manga, err := l.c.GetMangaInfo(ctx, mangaURL)
	if err != nil {
		return manga, err
	}
	manga.URL = theme.GetTrucattedURL(manga.URL)

	if len(manga.Chapters) == 0 {
		c := l.conn
		opts := theme.GetMadaraScrapeOptsForChapterList(c, manga.OtherID, mangaURL+"ajax/chapters/")
		chaptersManga, err := scrapper.ScrapeMangaInfo(ctx, c, &opts)
		if err != nil {
			return manga, err
		}

		for _, c := range chaptersManga.Chapters {
			c.URL = theme.GetTrucattedURL(c.URL)
			manga.Chapters = append(manga.Chapters, c)
		}
	}

	return manga, err
}

func (l *leviatan) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	chapterURL = theme.GetCompleteURL(chapterURL, l.GetSource().Domain)
	return l.c.GetChapterPages(ctx, chapterURL)
}
