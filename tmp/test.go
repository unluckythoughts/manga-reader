package main

import (
	"fmt"
	"net/http"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/PuerkitoBio/goquery"
	"github.com/unluckythoughts/manga-reader/sources"

	colly "github.com/gocolly/colly/v2"
)

func test1(url string) {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(100),
	)

	c.WithTransport(cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport))

	// Find and visit all links
	// c.OnHTML("div.listupd div.bs", func(e *colly.HTMLElement) {
	// 	fmt.Println("hasdf")
	// 	fmt.Println(e.DOM.Find("a[href]").Attr("href"))
	// })

	c.OnHTML("div#content", func(h *colly.HTMLElement) {
		val := h.DOM.Find("ul.clstyle li a span.chapterdate")

		val.Each(func(i int, s *goquery.Selection) {
			link, e := s.Attr("src")
			fmt.Printf("adsafasdf: %+v, \n %+v\n", strings.TrimSpace(link), e)
		})
	})

	c.OnResponseHeaders(func(r *colly.Response) {
		fmt.Printf("%+v\n", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %v failed with response: %+v\n Error: %+v", r.Request.URL, r.Ctx, err)
	})

	c.Visit(url)
}

func main() {

	url := "https://www.asurascans.com/reaper-of-the-drifting-moon-chapter-23/"

	// test1(url)

	a := sources.GetAsuraScansSource()
	manga, err := a.GetChapterImageURLs(url)
	fmt.Println(manga, err)

}
