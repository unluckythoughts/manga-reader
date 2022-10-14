package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetMangaClashConnector() models.IMangaConnector {
	c := theme.GetMadaraConnector()
	c.Source = models.Source{
		Name:    "Manga Clash",
		Domain:  "mangaclash.com",
		IconURL: "https://mangaclash.com/wp-content/uploads/2020/03/cropped-22.jpg",
	}
	c.BaseURL = "https://mangaclash.com/"
	c.MangaSelectors.Info.ChapterUploadDateFormat = "01/02/2006"

	return c
}
