package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetElarcPageConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Elarc Page",
		Domain:  "elarcpage.com",
		IconURL: "/static/assets/images/elarcpage.png",
	}
	c.BaseURL = "http://elarcpage.com/"
	c.MangaListPath = "series/"

	c.MangaSelectors.Info.Title = ".info-right h1.entry-title"
	c.MangaSelectors.Info.ImageURL = ".info-left .thumb img[src]"
	c.MangaSelectors.Info.Synopsis = ".info-right .wd-full .entry-content p"

	c.MangaSelectors.Chapter.ImageUrl = "#readerarea img[src]"

	return c
}
