package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetInfernalVoidScansConnector() models.IConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Invefernal Void Scans",
		Domain:  "void-scans.com",
		IconURL: "https://void-scans.com/wp-content/uploads/2021/09/cropped-weblogo-1.png",
	}
	c.BaseURL = "https://void-scans.com/"

	c.Selectors.Info.Title = ".info-right h1.entry-title"
	c.Selectors.Info.ImageURL = ".info-left .thumb img[src]"
	c.Selectors.Info.Synopsis = ".info-right .wd-full .entry-content p"

	c.Selectors.Chapter.ImageUrl = "#readerarea img[src]"

	return c
}
