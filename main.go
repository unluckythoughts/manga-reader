package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/book-reader/server/api"
	"github.com/unluckythoughts/go-microservice/v2"
	_ "github.com/unluckythoughts/go-microservice/v2/utils" // Import to register custom validators
	"go.uber.org/zap"
)

//go:embed public
var embeddedFiles embed.FS

// spaFileSystem wraps an http.FileSystem so that any missing path falls back
// to index.html, enabling client-side routing for the SPA.
type spaFileSystem struct{ root http.FileSystem }

func (fs spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open("/index.html")
	}
	return f, err
}

func proxyTransport(l *zap.Logger) http.RoundTripper {
	return cloudflarebp.AddCloudFlareByPass(http.DefaultTransport)
}

func main() {
	readFlags()
	opts := microservice.Options{
		Name:           "book-reader",
		EnableDB:       true,
		DBType:         microservice.DBTypeSqlite,
		ProxyTransport: proxyTransport,
	}
	s := microservice.New(opts)
	api.Register(s.HttpRouter(), s.GetDB(), s.GetLogger(), s.GetWorker())
	subFS, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		s.GetLogger().Fatal("Failed to create sub filesystem", zap.Error(err))
	}
	s.HttpRouter().ServeFiles("/app/*filepath", spaFileSystem{http.FS(subFS)})
	s.Start()
}
