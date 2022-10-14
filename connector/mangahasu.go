package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

func GetMangaHasuConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Manga Hasu",
		Domain:  "mangahasu.se",
		IconURL: "https://mangahasu.se/favicon.ico",
	}
	c.BaseURL = "http://mangahasu.se/"
	c.MangaListPath = "directory.html"

	c.List.MangaContainer = ".list_manga .div_item"
	c.List.MangaTitle = ".info-manga a h3"
	c.List.MangaURL = ".info-manga a[href]"
	c.List.MangaImageURL = ".wrapper_imagage img[src],img[src],.wrapper_imagage a[src]"
	c.List.LastPage = ".pagination-ct a[title='Trang cuối']"
	c.List.PageParam = "?page=" + scrapper.PAGE_ID

	c.Info.Title = ".wrapper_content .info-title h1"
	c.Info.ImageURL = ".wrapper_content .info-img img[src]"
	c.Info.Synopsis = ".wrapper_content .content-info > h3 + div"
	c.Info.ChapterContainer = ".wrapper_content .list-chapter tbody tr"
	c.Info.ChapterNumber = "td.name a"
	c.Info.ChapterTitle = "td.name a"
	c.Info.ChapterURL = "td.name a[href]"
	c.Info.ChapterUploadDate = "td.date-updated"
	c.Info.ChapterUploadDateFormat = "Jan 02, 2006"

	c.Chapter.ImageUrl = "#loadchapter .img img[data-src], #loadchapter .img img[src]"

	return c
}
