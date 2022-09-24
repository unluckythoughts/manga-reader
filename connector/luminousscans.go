package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetLuminousScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Luminous Scans",
		Domain:  "luminousscans.com",
		IconURL: "https://luminousscans.com/wp-content/uploads/2021/12/logo.png",
	}
	c.BaseURL = "http://luminousscans.com/"
	c.MangaListPath = "series/"
	c.MangaSelectors.Info.Title = ".info-right h1.entry-title"
	c.MangaSelectors.Info.ImageURL = ".info-left .thumb img[src]"
	c.MangaSelectors.Info.Synopsis = ".info-right .wd-full .entry-content"

	c.MangaSelectors.Chapter.ImageUrl = "#readerarea img[class*='wp-image'][src]"

	return c
}
