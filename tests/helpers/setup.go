package helpers

import (
	"github.com/unluckythoughts/manga-reader/tests/helpers/client"
	"github.com/unluckythoughts/manga-reader/tests/helpers/repo"
)

func Setup() {
	repo.Setup()
	client.Setup()
}
