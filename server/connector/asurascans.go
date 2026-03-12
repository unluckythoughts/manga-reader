package connector

import (
	"github.com/unluckythoughts/book-reader/server/connector/theme"
	"github.com/unluckythoughts/book-reader/server/models"
	"go.uber.org/zap"
)

type asura struct {
	*theme.BasicConnector
}

func NewAsuraConnector(l *zap.Logger) models.IConnector {
	conn := models.Connector{
		Name:        "AsuraScans",
		Domain:      "asurascanz.com",
		IconURL:     "https://asura-scans.online/wp-content/uploads/2025/05/logo.webp",
		Type:        models.BookTypeManga,
		BookListURL: "/manga/?status=&type=&order=title",
		Selectors: models.Selectors{
			BookListItem: "#content .listupd .bs",
			NextPage:     "#content .hpage a.r[href]",
			Book: models.BookSelectors{
				URL:             ".bsx a[href]",
				Title:           "h1.entry-title||.bsx a .bigor .tt",
				CoverImage:      ".main-info .thumb img[src]||.bsx a .limit img[src]",
				Synopsis:        ".main-info .info-right .entry-content p",
				ChapterListItem: "#chapterlist ul.clstyle li",
				Chapter: models.ChapterSelectors{
					URL:        ".main .top h1.tit a[href]||a[href]",
					Number:     "h1.entry-title||a span.chapternum",
					Title:      "h1.entry-title||a span.chapternum",
					UploadDate: "a span.chapterdate",
					DateFormat: "January 2, 2006",
					Content: models.ContentSelectors{
						Data: "#readerarea p img[src]",
					},
				},
			},
		},
	}

	bc := theme.NewBasic(conn, l).(*theme.BasicConnector)
	return &asura{bc}
}
