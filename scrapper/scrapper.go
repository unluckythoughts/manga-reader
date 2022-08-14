package scrapper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/unluckythoughts/manga-reader/models"
)

type Scrapper struct {
	src models.Source
}

func (s *Scrapper) getColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.WithTransport(s.src.RoundTripper)

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting URL: %+v\n", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		if err != nil {
			fmt.Printf("error requesting page %s, error %+v\n", r.Request.URL.String(), err)
		}
	})

	return c
}

func NewScrapper(src models.Source) *Scrapper {
	return &Scrapper{
		src: src,
	}
}
