package realmscans

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
	conn := connector.GetRealmScansConnector()
	return conn.GetMangaList(ctx)
}

func TestMangaList(t *testing.T) {
	// t.Skip("skipping to not download always")

	mangas, err := GetMangaList(t)
	assert.NoError(t, err, "error while getting asura scans mangalist")
	assert.Greater(t, len(mangas), 40, "could not get all the asura scans mangas")
}
