package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetInfernalVoidScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Infernal Void Scans",
		Domain:  "void-scans.com",
		IconURL: "/assets/images/voidscans.png",
	}
	c.BaseURL = "https://void-scans.com/"
	c.Chapter.ImageUrl = "#readerarea p img[data-lazy-src],#readerarea p img[src]"

	return c
}
