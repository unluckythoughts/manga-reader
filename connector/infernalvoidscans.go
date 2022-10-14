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
		IconURL: "https://void-scans.com/wp-content/uploads/cropped-weblogo-1.png",
	}
	c.BaseURL = "https://void-scans.com/"

	return c
}
