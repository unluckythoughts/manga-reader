package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

type reaper models.Connector

func GetReaperScansConnector() models.IConnector {
	c := theme.GetMadaraConnector()

	c.Source = models.Source{
		Name:    "Reaper Scans",
		Domain:  "reaperscans.com",
		IconURL: "https://reaperscans.com/wp-content/uploads/2021/07/cropped-ms-icon-310x310-1.png",
	}
	c.BaseURL = "https://reaperscans.com/"
	c.Selectors.Info.Title = ".post-title h1"
	c.Selectors.Info.ImageURL = ".tab-summary a img[data-src],.tab-summary a img[src],a#roi img[data-src],a#roi img[src],a#roiroi img[data-src],a#roiroi img[src]"
	c.Selectors.Info.Synopsis = ".container .summary__content"
	c.Selectors.Info.ChapterContainer = "ul.main li"
	c.Selectors.Info.ChapterNumber = ".chapter-link a > p"
	c.Selectors.Info.ChapterTitle = ".chapter-link a > p"
	c.Selectors.Info.ChapterURL = ".chapter-link a[href]"
	c.Selectors.Info.ChapterUploadDate = ".chapter-link a span i"
	c.Selectors.Info.ChapterUploadDateFormat = "Jan 02, 2006"

	return c
}
