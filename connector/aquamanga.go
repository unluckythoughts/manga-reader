package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

type aquamanga models.Connector

func GetAquaMangaConnector() models.IConnector {
	c := theme.GetMadaraConnector()
	c.Source = models.Source{
		Name:    "Aqua Manga",
		Domain:  "aquamanga.com",
		IconURL: "https://aquamanga.com/wp-content/uploads/2021/03/cropped-cropped-favicon-1.png",
	}
	c.BaseURL = "https://aquamanga.com/"

	return c
}
