package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetMangaHasuConnector() models.IConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Manga Hasu",
		Domain:  "mangahasu.se",
		IconURL: "https://mangahasu.se/favicon.ico",
	}
	c.BaseURL = "http://mangahasu.se/"
	c.MangaListPath = "directory.html"

	c.Selectors.List.MangaContainer = ".list_manga .div_item"
	c.Selectors.List.MangaTitle = ".info-manga a h3"
	c.Selectors.List.MangaURL = ".info-manga a[href]"
	c.Selectors.List.MangaImageURL = ".wrapper_imagage img[src],img[src],.wrapper_imagage a[src]"
	c.Selectors.List.NextPage = ".pagination-ct a[title='Tiếp']"
	c.Selectors.List.LastPage = ".pagination-ct a[title='Trang cuối']"

	c.Selectors.Info.Title = ".wrapper_content .info-title h1"
	c.Selectors.Info.ImageURL = ".wrapper_content .info-img img[src]"
	c.Selectors.Info.Synopsis = ".wrapper_content .content-info > h3 + div"
	c.Selectors.Info.ChapterContainer = ".wrapper_content .list-chapter tbody tr"
	c.Selectors.Info.ChapterNumber = "td.name a"
	c.Selectors.Info.ChapterTitle = "td.name a"
	c.Selectors.Info.ChapterURL = "td.name a[href]"
	c.Selectors.Info.ChapterUploadDate = "td.date-updated"
	c.Selectors.Info.ChapterUploadDateFormat = "Jan 02, 2006"

	c.Selectors.Chapter.ImageUrl = "#loadchapter .img img[data-src], #loadchapter .img img[src]"

	return c
}
