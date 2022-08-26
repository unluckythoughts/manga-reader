package leviatanscans

import (
	"os"
	"testing"

	"github.com/unluckythoughts/manga-reader/tests/helpers"
)

func cleanup() {}

func TestMain(m *testing.M) {
	helpers.Setup()
	code := m.Run()
	cleanup()
	os.Exit(code)
}
