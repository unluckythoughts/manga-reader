package main

import (
	"net/http"

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
	opts := microservice.Options{
		Name:           "manga-reader",
		EnableDB:       true,
		DBType:         microservice.DBTypeSqlite,
		ProxyTransport: proxyTransport,
	}
	s := microservice.New(opts)
	readerService := service.New(s.GetDB())

	reader.RegisterRoutes(s.HttpRouter(), readerService)
	s.HttpRouter().ServeFiles("/static/*filepath", http.Dir("./public"))
	s.Start()
}
