package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetFlameScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Flame Scans",
		Domain:  "flamescans.org",
		IconURL: "https://flamescans.org/favicon.ico",
	}
	c.BaseURL = "http://flamescans.org/"
	c.MangaListPath = "series/"
	c.MangaSelectors.Info.Title = ".info-half h1.entry-title"
	c.MangaSelectors.Info.ImageURL = ".thumb-half .thumb img[src]"
	c.MangaSelectors.Info.Synopsis = ".info-half .summary .wd-full .entry-content"

	return c
}
