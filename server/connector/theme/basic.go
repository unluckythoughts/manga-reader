package theme

import (
	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/book-reader/server/utils"
	"github.com/unluckythoughts/go-scraper"
)

type BasicConnector struct {
	s    *scraper.Scraper
	conn models.Connector
}

func NewBasic(conn models.Connector) models.IConnector {
	return &BasicConnector{
		s: scraper.New(scraper.Options{
			MaxRetries:          5,
			MaxParallelRequests: 4,
		}),
		conn: conn,
	}
}

func (c *BasicConnector) GetName() string {
	return c.conn.Name
}

func (c *BasicConnector) GetDomain() string {
	return c.conn.Domain
}

func (c *BasicConnector) GetIconURL() string {
	return c.conn.IconURL
}

func (c *BasicConnector) GetSelectors() models.Selectors {
	return c.conn.Selectors
}

func (c *BasicConnector) getCompleteURL(path string) string {
	baseURL := utils.GetBaseURL(c.GetDomain())
	return utils.GetCompleteURL(baseURL, path)
}

func (c *BasicConnector) GetBooks() ([]models.Book, error) {
	url := c.getCompleteURL(c.conn.BookListURL)
	config := scraper.PaginationConfig{
		NextPageSelector:   c.conn.Selectors.NextPage,
		LastPageSelector:   c.conn.Selectors.LastPage,
		NextPageURLPattern: c.conn.Selectors.NextPageURLPattern,
	}
	bookItemsChan, err := c.s.ScrapePaginated(url, c.conn.Selectors.BookListItem, config)
	if err != nil {
		return nil, err
	}

	books := []models.Book{}
	for bookItem := range bookItemsChan {
		if bookItem.Err != nil {
			return nil, bookItem.Err
		}

		book, err := getBook(bookItem.Data, c.conn)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (c *BasicConnector) GetBookChapters(bookURL string) ([]models.Chapter, error) {
	url := c.getCompleteURL(bookURL)
	config := scraper.PaginationConfig{
		NextPageSelector: c.conn.Selectors.Book.NextPage,
	}
	chapterItemsChan, err := c.s.ScrapePaginated(url, c.conn.Selectors.Book.ChapterListItem, config)
	if err != nil {
		return nil, err
	}
	chapters := []models.Chapter{}
	for chapterItem := range chapterItemsChan {
		chapter, err := getChapter(chapterItem.Data, c.conn)
		if err != nil {
			return nil, err
		}
		chapters = append(chapters, chapter)
	}
	return chapters, nil
}

func (c *BasicConnector) GetChapterContent(chapterURL string) (models.List, error) {
	url := c.getCompleteURL(chapterURL)
	var content models.List

	htmlContent, err := c.s.ScrapeHTML(url)
	if err != nil {
		return content, err
	}

	contentData, err := scraper.GetText(htmlContent, c.conn.Selectors.Book.Chapter.Content.Data)
	if err != nil {
		return content, err
	}

	patterns := c.conn.Selectors.Book.Chapter.Content.ReplacePatterns
	for _, data := range contentData {
		switch c.conn.Type {
		case models.BookTypeNovel:
			content.Add(cleanData(data, patterns))
		case models.BookTypeManga:
			content.Add(data)
		}
	}

	return content, nil
}
