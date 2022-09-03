package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetLeviatanScansConnector() models.IConnector {
	c := theme.GetMadaraConnector()

	c.Source = models.Source{
		Name:    "Leviatan Scans",
		Domain:  "leviatanscans.com",
		IconURL: "https://leviatanscans.com/wp-content/uploads/2021/03/cropped-isotiponegro.png",
	}

	c.BaseURL = "https://leviatanscans.com/"
	c.Selectors.Info.Title = "#manga-title h1"
	c.Selectors.Info.OtherID = "#manga-chapters-holder[data-id]"
	c.Selectors.Info.Synopsis = ".summary_content .post-content_item:last-of-type p span"
	c.Selectors.Info.ChapterUploadDateFormat = "Jan 2, 2006"

	return c
}
