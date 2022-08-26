package mangahub

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
	conn := connector.GetMangaHubConnector()
	manga, err := conn.GetMangaInfo(ctx, mangaUrl)
	assert.NoError(t, err, "error while getting asura scans manga")
	assert.Equal(t, manga.URL, mangaUrl, "manga url mismatch")
	assert.Regexp(t, `"Aoi" - Hikaru ga Chikyuu ni Itakoro......`, manga.Title, "manga title mismatch")
	assert.NotEmpty(t, manga.ImageURL, "manga image url is empty")
	assert.NotEmpty(t, manga.Synopsis, "manga synopsis is empty")
	assert.Greater(t, len(manga.Chapters), 14, "could not all the chapters for manga")

	return manga, err
}

func TestManga(t *testing.T) {
	t.Skip("skipping to not download always")

	url := "https://mangahub.io/manga/aoi-hikaru-ga-chikyuu-ni-itakoro"
	GetMangaInfo(t, url)
}
