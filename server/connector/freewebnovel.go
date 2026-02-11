package connector

import (
	"github.com/unluckythoughts/book-reader/server/connector/theme"
	"github.com/unluckythoughts/book-reader/server/models"
	"go.uber.org/zap"
)

type fwb struct {
	*theme.BasicConnector
}

func NewFreeWebNovelConnector(l *zap.Logger) models.IConnector {
	conn := models.Connector{
		Name:        "FreeWebNovel",
		Domain:      "freewebnovel.com",
		IconURL:     "/static/freewebnovel/favicon.ico",
		Type:        models.BookTypeNovel,
		BookListURL: "/sort/latest-novel/",
		Selectors: models.Selectors{
			BookListItem:       ".main .wp .ul-list1 .li .con",
			LastPage:           ".main .wp .pages ul li a:last-child[href]",
			NextPageURLPattern: "/sort/latest-novel/::page::",
			Book: models.BookSelectors{
				URL:             ".txt h3.tit a[href]",
				Title:           ".main .wp .m-book1 .m-desc h1.tit||.txt h3.tit a",
				CoverImage:      ".main .wp .m-book1 .m-imgtxt .pic img[src]||.pic a img[src]",
				Synopsis:        ".m-desc .txt .inner p",
				ChapterListItem: "ul#idData li",
				Chapter: models.ChapterSelectors{
					URL:    ".main .top h1.tit a[href]||a[href]",
					Number: ".main .top h1.tit a||a",
					Title:  ".main .top h1.tit a||a",
					Content: models.ContentSelectors{
						Data:            "#article p",
						ReplacePatterns: []models.Pattern{},
					},
				},
			},
		},
	}

	bc := theme.NewBasic(conn, l).(*theme.BasicConnector)
	return &fwb{bc}
}
