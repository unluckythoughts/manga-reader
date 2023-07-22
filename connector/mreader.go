package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/scrapper"

	"github.com/unluckythoughts/manga-reader/models"
)

func GetMReaderConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "MReader",
		Domain:  "mreader.co",
		IconURL: "https://mreader.co/static/img/logo_200x200.png",
	}
	c.BaseURL = "http://mreader.co/"
	c.MangaListPath = "browse-comics/?filter=Updated"

	c.List.MangaContainer = "ul.novel-list li.novel-item"
	c.List.MangaTitle = "a h4.novel-title"
	c.List.MangaURL = "a[href]"
	c.List.MangaImageURL = "a .cover-wrap img[data-src], a .cover-wrap img[src]"
	c.List.MangaSlug = ""
	c.List.MangaOtherID = ""
	c.List.LastPage = "ul.pagination li:nth-of-type(2)"
	c.List.PageParam = "&results=" + scrapper.PAGE_ID

	c.Info.Title = ".main-head h1.novel-title"
	c.Info.ImageURL = "figure.cover img[src]"
	c.Info.Synopsis = "section#info p.description"
	c.Info.Slug = ""
	c.Info.OtherID = ""
	c.Info.ChapterListURL = "#chpagedlist .intro a[href]"
	c.Info.ChapterListNextPage = ""
	c.Info.ChapterListLastPage = ""
	c.Info.ChapterListPageParam = ""
	c.Info.ChapterContainer = "section#chapters ul.chapter-list li"
	c.Info.ChapterNumber = "a span.chapter-no, a strong.chapter-title"
	c.Info.ChapterTitle = "a strong.chapter-title"
	c.Info.ChapterURL = "a[href]"
	c.Info.ChapterUploadDate = "time.chapter-update[datetime]"
	c.Info.ChapterUploadDateFormat = "Jan. 02,2006, 3:10 a.m."

	c.Chapter.ImageUrl = "#chapter-reader img[src]"

	return c
}
