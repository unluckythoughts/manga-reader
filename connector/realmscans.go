package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetRealmScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Realm Scans",
		Domain:  "realmscans.xyz",
		IconURL: "/static/assets/images/realmscans.svg",
	}
	c.BaseURL = "https://realmscans.xyz/"
	c.MangaListPath = "series/"

	c.MangaSelectors.Info.Title = ".info-right h1"
	c.MangaSelectors.Info.ImageURL = ".thumb img[src]"
	c.MangaSelectors.Info.Synopsis = ".info-right .wd-full .entry-content p"

	c.MangaSelectors.Chapter.ImageUrl = "#readerarea img:is([data-src],[src])"
	c.Info.ChapterUploadDateFormat = "02 Jan 2006"

	return c
}
