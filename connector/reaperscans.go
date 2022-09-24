package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

type reaper models.MangaConnector

func GetReaperScansConnector() models.IMangaConnector {
	c := theme.GetMadaraConnector()

	c.Source = models.Source{
		Name:    "Reaper Scans",
		Domain:  "reaperscans.com",
		IconURL: "https://reaperscans.com/wp-content/uploads/2021/07/cropped-ms-icon-310x310-1.png",
	}
	c.BaseURL = "https://reaperscans.com/"
	c.MangaSelectors.Info.Title = ".post-title h1"
	c.MangaSelectors.Info.ImageURL = ".tab-summary a img[data-src],.tab-summary a img[src],a#roi img[data-src],a#roi img[src],a#roiroi img[data-src],a#roiroi img[src]"
	c.MangaSelectors.Info.Synopsis = ".container .summary__content"
	c.MangaSelectors.Info.ChapterContainer = "ul.main li"
	c.MangaSelectors.Info.ChapterNumber = ".chapter-link a > p"
	c.MangaSelectors.Info.ChapterTitle = ".chapter-link a > p"
	c.MangaSelectors.Info.ChapterURL = ".chapter-link a[href]"
	c.MangaSelectors.Info.ChapterUploadDate = ".chapter-link a span i"
	c.MangaSelectors.Info.ChapterUploadDateFormat = "Jan 02, 2006"

	return c
}
