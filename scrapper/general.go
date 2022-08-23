package scrapper

import (
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
)

func GetHeaders(ctx web.Context, url string, roundTripper http.RoundTripper) http.Header {
	c := getColly(ctx, roundTripper)

	var respHeaders http.Header
	c.OnResponseHeaders(func(r *colly.Response) {
		respHeaders = (*r).Headers.Clone()
	})

	c.Visit(url)
	return respHeaders
}
