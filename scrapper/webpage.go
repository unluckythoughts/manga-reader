package scrapper

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"go.uber.org/zap"
)

func GetPageForScrapping(ctx web.Context, opts *ScrapeOptions, cb colly.HTMLCallback) error {
	opts.SetDefaults()
	c := getColly(ctx, opts.RoundTripper)

	if opts.InitialHtmlTag == WHOLE_BODY_TAG {
		c.OnResponse(func(r *colly.Response) {
			dom, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debug("error reading response body")
				return
			}

			i := 0
			for _, n := range dom.Selection.Nodes {
				el := colly.NewHTMLElementFromSelectionNode(r, dom.Selection, n, i)
				i++
				cb(el)
			}
		})
	} else {
		c.OnHTML(opts.InitialHtmlTag, cb)
	}

	err := c.Request(opts.RequestMethod, opts.URL, opts.Body, nil, opts.Headers)
	if err != nil {
		return err
	}

	return nil
}
