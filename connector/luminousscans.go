package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetLuminousScansConnector() models.IConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Luminous Scans",
		Domain:  "luminousscans.com",
		IconURL: "https://luminousscans.com/wp-content/uploads/2021/12/logo.png",
	}
	c.BaseURL = "http://luminousscans.com/"
	c.MangaListPath = "series/"
	c.Selectors.Info.Title = ".info-right h1.entry-title"
	c.Selectors.Info.ImageURL = ".info-left .thumb img[src]"
	c.Selectors.Info.Synopsis = ".info-right .summary .wd-full .entry-content"

	return c
}
