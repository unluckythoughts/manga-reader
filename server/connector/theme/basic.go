package theme

import (
	"fmt"

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
	return "", fmt.Errorf("To be implemented")
}
