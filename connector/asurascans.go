package connector

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector/theme"

	"github.com/unluckythoughts/manga-reader/models"
)

type asura struct {
	c    models.IMangaConnector
	conn models.MangaConnector
}

func GetAsuraScansConnector() models.IMangaConnector {
	c := theme.GetBasicWordPressConnector()

	c.Source = models.Source{
		Name:    "Asura Scans",
		Domain:  "asura.gg",
		IconURL: "/assets/images/asurascans.png",
	}
	c.BaseURL = "http://asura.gg/"
	c.Chapter.ImageUrl = "#readerarea img[src]"

	return &asura{c: c, conn: models.MangaConnector(*c)}
}

func (a *asura) GetSource() models.Source {
	return a.c.GetSource()
}

func (a *asura) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	return a.c.GetMangaList(ctx)
}

func (a *asura) GetLatestMangaList(ctx web.Context, latestTitle string) ([]models.Manga, error) {
	return a.c.GetLatestMangaList(ctx, latestTitle)
}

func (a *asura) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	return a.c.GetMangaInfo(ctx, mangaURL)
}

func (a *asura) GetChapterPages(ctx web.Context, pageListURL string) (models.Pages, error) {
	return a.c.GetChapterPages(ctx, pageListURL)
}
