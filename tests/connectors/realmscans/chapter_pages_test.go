package realmscans

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unluckythoughts/go-microservice/tools/logger"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetChapterPages(t *testing.T, chapterURL string) (models.Pages, error) {
	t.Helper()

	ctx := web.NewContext(logger.New(logger.Options{}))
	conn := connector.GetRealmScansConnector()
	pages, err := conn.GetChapterPages(ctx, chapterURL)
	assert.NoError(t, err, "error while getting asura scans chapter images")
	assert.Greater(t, len(pages.URLs), 32, "could not all the pages for chapter")
	for _, imageURL := range pages.URLs {
		assert.NotEmpty(t, imageURL, "could not get image url for chapter")

	}

	return pages, err
}

func TestChapterPages(t *testing.T) {
	t.Skip("skipping to not download always")

	url := "https://realmscans.com/infinite-level-up-murim-chapter-73/"
	GetChapterPages(t, url)
}
