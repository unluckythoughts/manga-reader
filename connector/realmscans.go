package connector

import (
	"github.com/unluckythoughts/manga-reader/connector/theme"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetRealmScansConnector() models.IConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Realm Scans",
		Domain:  "realmscans.com",
		IconURL: "https://realmscans.com/wp-content/uploads/2022/08/realm-scans-fav.png",
	}
	c.BaseURL = "http://realmscans.com/"
	c.MangaListPath = "series/"

	c.Selectors.Info.Title = ".info-right h1"
	c.Selectors.Info.ImageURL = ".thumb img[src]"
	c.Selectors.Info.Synopsis = ".info-right .wd-full .entry-content p"

	c.Selectors.Chapter.ImageUrl = "#readerarea img:is([data-src],[src])"

	return c
}
