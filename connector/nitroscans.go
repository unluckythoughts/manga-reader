package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetNitroScansConnector() models.IMangaConnector {
	c := theme.GetMadaraConnector()

	c.Source = models.Source{
		Name:    "Nitro Scans",
		Domain:  "nitroscans.com",
		IconURL: "https://nitroscans.com/wp-content/uploads/2021/05/cropped-132x132-1.png",
	}
	c.BaseURL = "https://nitroscans.com/"
	c.Info.ChapterListURL = ""

	return c
}
