package scrapper

import (
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

type Scrapper struct {
	src models.Source
}

func GetHeaders(ctx web.Context, url string, roundTripper http.RoundTripper) http.Header {
	c := getColly(ctx, roundTripper)

	var respHeaders http.Header
	c.OnResponseHeaders(func(r *colly.Response) {
		respHeaders = (*r).Headers.Clone()
	})

	c.Visit(url)
	return respHeaders
}

func getColly(ctx web.Context, rt http.RoundTripper) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	if rt != nil {
		c.WithTransport(rt)
	}

	c.OnRequest(func(r *colly.Request) {
		ctx.Logger().Debugf("Visiting URL: %+v", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		if err != nil {
			ctx.Logger().Debugf("error requesting page %s, error %+v", r.Request.URL.String(), err)
		}
	})

	return c
}
