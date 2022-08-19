package main

import (
	"net/http"
	"os"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice"
	"github.com/unluckythoughts/manga-reader/reader"
	"github.com/unluckythoughts/manga-reader/reader/service"
	"go.uber.org/zap"
)

func proxyTransport(l *zap.Logger) http.RoundTripper {
	return cloudflarebp.AddCloudFlareByPass(http.DefaultTransport)
}

func main() {
	_ = os.Setenv("DB_FILE_PATH", "db.sqlite")
	_ = os.Setenv("WEB_PORT", "5678")
	_ = os.Setenv("WEB_CORS", "true")
	_ = os.Setenv("WEB_PROXY", "true")
	_ = os.Setenv("DB_DEBUG", "true")

	opts := microservice.Options{
		Name:           "manga-reader",
		EnableDB:       true,
		DBType:         microservice.DBTypeSqlite,
		ProxyTransport: proxyTransport,
	}
	s := microservice.New(opts)

	readerService := service.New(s.GetDB())
	reader.RegisterRoutes(s.HttpRouter(), readerService)
	s.Start()
}
