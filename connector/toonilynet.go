package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetToonilyNetConnector() models.IMangaConnector {
	c := theme.GetMadaraConnector()
	c.Source = models.Source{
		Name:    "Toonily Net",
		Domain:  "toonily.net",
		IconURL: "https://toonily.net/wp-content/uploads/2020/06/cropped-android-chrome-512x512-1-32x32.png",
	}
	c.BaseURL = "https://toonily.net/"
	c.MangaSelectors.Info.ChapterUploadDateFormat = "January 2, 2006"

	return c
}
