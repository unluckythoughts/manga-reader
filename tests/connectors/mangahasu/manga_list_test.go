package mangahasu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unluckythoughts/go-microservice/tools/logger"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetMangaList(t *testing.T) ([]models.Manga, error) {
	t.Helper()

	ctx := web.NewContext(logger.New(logger.Options{}))
	conn := connector.GetMangaHasuConnector()

	mangas, err := conn.GetMangaList(ctx)
	assert.NoError(t, err, "error while getting asura scans mangalist")
	assert.Greater(t, len(mangas), 100, "could not get all the asura scans mangas")

	return mangas, err
}

func TestMangaList(t *testing.T) {
	t.Skip("skipping to not download always")

	GetMangaList(t)
}
