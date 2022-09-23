package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"

	"github.com/unluckythoughts/manga-reader/models"
)

func GetAsuraScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Asura Scans",
		Domain:  "asurascans.com",
		IconURL: "https://www.asurascans.com/wp-content/uploads/2021/03/Group_1.png",
	}
	c.BaseURL = "http://asurascans.com/"

	return c
}
