package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetMangaReadConnector() models.IMangaConnector {
	c := theme.GetMadaraConnector()

	c.Source = models.Source{
		Name:    "MangaRead",
		Domain:  "mangaread.org",
		IconURL: "/static/assets/images/mangaread.svg",
	}
	c.BaseURL = "https://www.mangaread.org/"
	c.Info.ChapterUploadDateFormat = "02.01.2006"

	return c
}
