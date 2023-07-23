package connector

import (
	"net/url"

	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
	"github.com/unluckythoughts/manga-reader/utils"
)

func GetReaperScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Reaper Scans",
		Domain:  "reaperscans.com",
		IconURL: "/static/assets/images/reaperscans.png",
	}
	c.BaseURL = "https://reaperscans.com/"
	c.MangaListPath = "comics/"
	c.MangaListURLParams = url.Values{}

	c.List.MangaContainer = "main li"
	c.List.MangaTitle = "a:last-of-type"
	c.List.MangaURL = "a:last-of-type[href]"
	c.List.MangaImageURL = "a img[src]"
	c.List.LastPage = "DEFAULT::6"
	c.List.PageParam = "?page=" + scrapper.PAGE_ID

	c.Info.Title = ".container h1"
	c.Info.ImageURL = ".container img[src]"
	c.Info.Synopsis = "section p.w-full"
	c.Info.ChapterContainer = ".container ~ div ul[role='list'] li"
	c.Info.ChapterListLastPage = "span:nth-last-child(2) > button[aria-label]"
	c.Info.ChapterListPageParam = "?page=" + scrapper.PAGE_ID
	c.Info.ChapterNumber = "a p.truncate"
	c.Info.ChapterTitle = "a p.truncate"
	c.Info.ChapterURL = "a[href]"
	// TODO: get all chapters by page
	c.Info.ChapterUploadDate = "a .mt-2 p"
	c.Info.ChapterUploadDateFormat = utils.HUMAN_READABLE_DATE_FORMAT

	c.Chapter.ImageUrl = "main img.display-block[src]"

	return c
}
