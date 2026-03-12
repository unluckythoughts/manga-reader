package theme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/book-reader/server/utils"
	"github.com/unluckythoughts/go-scraper"
	"go.uber.org/zap"
)

type BasicConnector struct {
	s    *scraper.Scraper
	l    *zap.Logger
	conn models.Connector
}

func NewBasic(conn models.Connector, logger *zap.Logger) models.IConnector {
	return &BasicConnector{
		s: scraper.New(scraper.Options{
			UserAgent:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			MaxRetries:          5,
			UseCloudflareBypass: true,
			MaxParallelRequests: 2, // Reduce parallelism to appear less bot-like
		}),
		conn: conn,
		l:    logger,
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

func (c *BasicConnector) SupportsLastPageSelector() bool {
	if c.conn.Selectors.LastPage == "" {
		return false
	}

	if c.conn.Selectors.NextPageURLPattern == "" {
		return false
	}

	if !strings.Contains(c.conn.Selectors.NextPageURLPattern, "::page::") {
		return false
	}

	return true
}

func (c *BasicConnector) GetBooksAsync() (<-chan models.Book, error) {
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

	booksChan := make(chan models.Book)

	go func() {
		defer close(booksChan)
		for bookItem := range bookItemsChan {
			if bookItem.Err != nil {
				c.l.Debug("Error while scraping book item", zap.Error(bookItem.Err))
				continue
			}

			book, err := getBook(bookItem.Data, c.conn)
			if err != nil {
				c.l.Debug("Error while getting book", zap.Error(err))
				continue
			}

			booksChan <- book
		}
	}()

	return booksChan, nil
}

func (c *BasicConnector) GetBookSynopsis(bookURL string) (string, error) {
	url := c.getCompleteURL(bookURL)
	html, err := c.s.ScrapeHTML(url)
	if err != nil {
		return "", err
	}

	texts, err := scraper.GetText(html, c.conn.Selectors.Book.Synopsis)
	if err != nil {
		return "", err
	}

	synopsis := strings.Join(texts, "\n")

	return synopsis, nil
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

func (c *BasicConnector) GetBookCount() (int, []models.Book, error) {
	if c.conn.Selectors.LastPage == "" {
		books, err := c.GetBooks()
		if err != nil {
			return 0, books, err
		}
		return len(books), books, nil
	}
	url := c.getCompleteURL(c.conn.BookListURL)
	html, err := c.s.ScrapeHTML(url)
	if err != nil {
		return 0, nil, err
	}

	bookItems, err := scraper.GetOuterHTML(html, c.conn.Selectors.BookListItem)
	if err != nil {
		return 0, nil, err
	}
	booksPerPageCount := len(bookItems)

	books := []models.Book{}
	for _, bookItem := range bookItems {
		book, err := getBook(bookItem, c.conn)
		if err != nil {
			return 0, nil, err
		}
		books = append(books, book)
	}

	lastPageItems, err := scraper.GetOuterHTML(html, c.conn.Selectors.LastPage)
	if err != nil {
		return 0, nil, err
	} else if len(lastPageItems) == 0 {
		return 0, nil, fmt.Errorf("error while getting books count: could not get lastpage number")
	}

	lastPageNum, err := scraper.GetInt(lastPageItems[0], "")
	if err != nil {
		return 0, nil, err
	}

	if lastPageNum == 0 {
		return 0, nil, fmt.Errorf("error while getting books count: last page number is zero")
	}

	lastPageURL := strings.ReplaceAll(
		c.conn.Selectors.NextPageURLPattern,
		"::page::", strconv.Itoa(lastPageNum),
	)
	lastPageURL = scraper.GetFullURL(url, lastPageURL)

	lastPageBookItems, err := c.s.ScrapeOuterHTML(lastPageURL, c.conn.Selectors.BookListItem)
	if err != nil {
		return 0, nil, err
	}
	lastPageBooksCount := len(lastPageBookItems)
	totalBooks := (lastPageNum-1)*booksPerPageCount + lastPageBooksCount

	return totalBooks, books, nil
}
