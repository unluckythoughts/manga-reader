package asurascans

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unluckythoughts/go-microservice/tools/logger"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetMangaInfo(t *testing.T, mangaUrl string) (models.Manga, error) {
	t.Helper()

	ctx := web.NewContext(logger.New(logger.Options{}))
	conn := connector.GetAsuraScansConnector()
	manga, err := conn.GetMangaInfo(ctx, mangaUrl)
	assert.NoError(t, err, "error while getting asura scans manga")
	assert.Equal(t, manga.URL, mangaUrl, "manga url mismatch")
	assert.Regexp(t, "To Hell With Being A Saint, I.m A Doctor", manga.Title, "manga title mismatch")
	assert.NotEmpty(t, manga.ImageURL, "manga image url is empty")
	assert.NotEmpty(t, manga.Synopsis, "manga synopsis is empty")
	assert.Greater(t, len(manga.Chapters), 31, "could not all the chapters for manga")

	return manga, err
}

func TestManga(t *testing.T) {
	t.Skip("skipping to not download always")

	url := "https://www.asurascans.com/manga/1660333069-to-hell-with-being-a-saint-im-a-doctor/"
	GetMangaInfo(t, url)
}
