package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetElarcPageConnector() models.IConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Elarc Page",
		Domain:  "elarcpage.com",
		IconURL: "https://elarcpage.com/favicon.ico",
	}
	c.BaseURL = "http://asurascans.com/"

	c.Selectors.Info.Title = ".info-right h1.entry-title"
	c.Selectors.Info.ImageURL = ".info-left .thumb img[src]"
	c.Selectors.Info.Synopsis = ".info-right .wd-full .entry-content p"

	c.Selectors.Chapter.ImageUrl = "#readerarea img[src]"

	return c
}
