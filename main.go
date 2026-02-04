package main

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/book-reader/server/api"
	"github.com/unluckythoughts/go-microservice/v2"
	"go.uber.org/zap"
)

func proxyTransport(l *zap.Logger) http.RoundTripper {
	return cloudflarebp.AddCloudFlareByPass(http.DefaultTransport)
}

func main() {
	opts := microservice.Options{
		Name:           "book-reader",
		EnableDB:       true,
		DBType:         microservice.DBTypeSqlite,
		ProxyTransport: proxyTransport,
	}
	s := microservice.New(opts)
	api.Register(s.HttpRouter(), s.GetDB(), s.GetLogger(), s.GetWorker())
	s.HttpRouter().ServeFiles("/static/*filepath", http.Dir("./public"))
	s.Start()
}
