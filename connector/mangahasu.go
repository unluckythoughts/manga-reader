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

	c.MangaSelectors.List.MangaContainer = ".list_manga .div_item"
	c.MangaSelectors.List.MangaTitle = ".info-manga a h3"
	c.MangaSelectors.List.MangaURL = ".info-manga a[href]"
	c.MangaSelectors.List.MangaImageURL = ".wrapper_imagage img[src],img[src],.wrapper_imagage a[src]"
	c.MangaSelectors.List.LastPage = ".pagination-ct a[title='Trang cuối']"
	c.MangaSelectors.List.PageParam = "?page=" + scrapper.MANGA_LIST_PAGE_ID

	c.MangaSelectors.Info.Title = ".wrapper_content .info-title h1"
	c.MangaSelectors.Info.ImageURL = ".wrapper_content .info-img img[src]"
	c.MangaSelectors.Info.Synopsis = ".wrapper_content .content-info > h3 + div"
	c.MangaSelectors.Info.ChapterContainer = ".wrapper_content .list-chapter tbody tr"
	c.MangaSelectors.Info.ChapterNumber = "td.name a"
	c.MangaSelectors.Info.ChapterTitle = "td.name a"
	c.MangaSelectors.Info.ChapterURL = "td.name a[href]"
	c.MangaSelectors.Info.ChapterUploadDate = "td.date-updated"
	c.MangaSelectors.Info.ChapterUploadDateFormat = "Jan 02, 2006"

	c.MangaSelectors.Chapter.ImageUrl = "#loadchapter .img img[data-src], #loadchapter .img img[src]"

	return c
}
