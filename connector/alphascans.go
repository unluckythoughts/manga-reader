package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"

	"github.com/unluckythoughts/manga-reader/models"
)

func GetAlphaScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Alpha Scans",
		Domain:  "alpha-scans.org",
		IconURL: "https://alpha-scans.org/wp-content/uploads/2022/02/website-alpha-logo-copy.png",
	}
	c.BaseURL = "https://alpha-scans.org/"
	c.MangaListPath = "manga/"

	return c
}
