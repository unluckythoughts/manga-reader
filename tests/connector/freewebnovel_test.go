package connector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/unluckythoughts/book-reader/server/connector"
	"github.com/unluckythoughts/book-reader/server/models"
)

type FreeWebNovelTestSuite struct {
	ConnectorTestSuite
	connector models.IConnector
}

func TestFreeWebNovelTestSuite(t *testing.T) {
	suite.Run(t, new(FreeWebNovelTestSuite))
}

func (s *FreeWebNovelTestSuite) SetupSuite() {
	s.connector = connector.NewFreeWebNovelConnector()
}

// isRodBlockedError checks if the error is due to Windows blocking rod/leakless
func (s *FreeWebNovelTestSuite) isRodBlockedError(err error) bool {
	if err == nil {
		return false
	}
	errorMsg := err.Error()
	// Check for Windows Defender or antivirus blocking messages
	return strings.Contains(errorMsg, "leakless.exe") ||
		strings.Contains(errorMsg, "virus") ||
		strings.Contains(errorMsg, "potentially unwanted software")
}

func (s *FreeWebNovelTestSuite) TestConnectorProperties() {
	assert.NotNil(s.T(), s.connector, "Connector should not be nil")
}

func (s *FreeWebNovelTestSuite) TestGetName() {
	name := s.connector.GetName()
	assert.Equal(s.T(), "FreeWebNovel", name, "Connector name should be 'FreeWebNovel'")
}

func (s *FreeWebNovelTestSuite) TestGetDomain() {
	domain := s.connector.GetDomain()
	assert.Equal(s.T(), "freewebnovel.com", domain, "Domain should be 'freewebnovel.com'")
}

func (s *FreeWebNovelTestSuite) TestGetIconURL() {
	iconURL := s.connector.GetIconURL()
	assert.Equal(s.T(), "/static/freewebnovel/favicon.ico", iconURL, "Icon URL should match expected path")
}

func (s *FreeWebNovelTestSuite) TestGetSelectors() {
	selectors := s.connector.GetSelectors()

	assert.NotEmpty(s.T(), selectors.BookListItem, "BookListItem selector should not be empty")
	assert.NotEmpty(s.T(), selectors.LastPage, "LastPage selector should not be empty")
	assert.NotEmpty(s.T(), selectors.NextPageURLPattern, "NextPageURLPattern should not be empty")

	// Test book selectors
	assert.NotEmpty(s.T(), selectors.Book.URL, "Book URL selector should not be empty")
	assert.NotEmpty(s.T(), selectors.Book.Title, "Book Title selector should not be empty")
	assert.NotEmpty(s.T(), selectors.Book.CoverImage, "Book CoverImage selector should not be empty")
	assert.NotEmpty(s.T(), selectors.Book.Synopsis, "Book Synopsis selector should not be empty")
	assert.NotEmpty(s.T(), selectors.Book.ChapterListItem, "ChapterListItem selector should not be empty")

	// Test chapter selectors
	assert.NotEmpty(s.T(), selectors.Book.Chapter.URL, "Chapter URL selector should not be empty")
	assert.NotEmpty(s.T(), selectors.Book.Chapter.Number, "Chapter Number selector should not be empty")
	assert.NotEmpty(s.T(), selectors.Book.Chapter.Title, "Chapter Title selector should not be empty")
	assert.NotEmpty(s.T(), selectors.Book.Chapter.Content.Data, "Chapter Content Data selector should not be empty")
}

func (s *FreeWebNovelTestSuite) TestGetBookCount() {
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}

	count, books, err := s.connector.GetBookCount()

	if err != nil {
		s.T().Logf("Warning: GetBookCount failed (site may be down or blocking requests): %v", err)
		if s.isRodBlockedError(err) {
			s.T().Skip("Skipping test: Windows is blocking browser automation (leakless.exe). See README for instructions.")
		}
		s.T().Skip("Skipping test due to network/site availability issues")
		return
	}

	assert.NoError(s.T(), err, "GetBookCount should not return an error")
	assert.Greater(s.T(), count, 0, "Book count should be greater than 0")

	// books may be nil if last page selector is present
	if books != nil {
		assert.NotEmpty(s.T(), books, "Books list should not be empty when returned")
	}
}

func (s *FreeWebNovelTestSuite) TestGetBookChapters() {
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}

	// First get a book to test with
	_, books, err := s.connector.GetBookCount()
	if err != nil {
		s.T().Skip("Skipping test: cannot get books for testing")
		return
	}

	if len(books) == 0 {
		s.T().Skip("Skipping test: no books available")
		return
	}

	bookURL := books[0].URL

	chapters, err := s.connector.GetBookChapters(bookURL, "")

	if err != nil {
		s.T().Logf("Warning: GetBookChapters failed: %v", err)
		s.T().Skip("Skipping test due to network/site availability issues")
		return
	}

	assert.NoError(s.T(), err, "GetBookChapters should not return an error")
	assert.NotEmpty(s.T(), chapters, "Chapters list should not be empty")

	// Validate first chapter structure
	if len(chapters) > 0 {
		chapter := chapters[0]
		assert.NotEmpty(s.T(), chapter.Title, "Chapter title should not be empty")
		assert.NotEmpty(s.T(), chapter.URL, "Chapter URL should not be empty")
		assert.NotEmpty(s.T(), chapter.Number, "Chapter number should not be empty")
	}
}

func (s *FreeWebNovelTestSuite) TestGetBookChaptersWithFilter() {
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}

	// First get a book to test with
	_, books, err := s.connector.GetBookCount()
	if err != nil {
		s.T().Skip("Skipping test: cannot get books for testing")
		return
	}

	if len(books) == 0 {
		s.T().Skip("Skipping test: no books available")
		return
	}

	bookURL := books[0].URL

	// Get all chapters first
	allChapters, err := s.connector.GetBookChapters(bookURL, "")
	if err != nil || len(allChapters) == 0 {
		s.T().Skip("Skipping test: cannot get chapters for testing")
		return
	}

	// Try to filter chapters greater than chapter 1
	filteredChapters, err := s.connector.GetBookChapters(bookURL, "1")
	assert.NoError(s.T(), err, "GetBookChapters on "+bookURL+" with filter should not return an error")

	// Filtered chapters should be less than or equal to all chapters
	assert.LessOrEqual(s.T(), len(filteredChapters), len(allChapters),
		"Filtered chapters should be less than or equal to all chapters")
}

func (s *FreeWebNovelTestSuite) TestGetChapterContent() {
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}

	// First get a book and then a chapter
	_, books, err := s.connector.GetBookCount()
	assert.NoError(s.T(), err, "GetBookCount should not return an error")
	assert.NotEmpty(s.T(), books, "Books list should not be empty")

	chapters, err := s.connector.GetBookChapters(books[0].URL, "")
	assert.NoError(s.T(), err, "GetBookChapters should not return an error")
	assert.NotEmpty(s.T(), chapters, "Chapters list should not be empty")

	chapterURL := chapters[0].URL

	content, err := s.connector.GetChapterContent(chapterURL)
	assert.NoError(s.T(), err, "GetChapterContent should not return an error")
	assert.NotEmpty(s.T(), content.Values(), "Chapter content should not be empty")
	assert.Greater(s.T(), content.Count(), 0, "Content should have at least one paragraph")

	// Validate content structure
	contentValues := content.Values()
	for _, paragraph := range contentValues {
		assert.NotEmpty(s.T(), paragraph, "Each content paragraph should not be empty")
	}
}

func (s *FreeWebNovelTestSuite) TestIntegrationFlow() {
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}

	// Test the complete flow: Books -> Chapters -> Content

	// Step 1: Get books
	_, books, err := s.connector.GetBookCount()
	if err != nil {
		s.T().Skip("Skipping integration flow test: cannot get books")
		return
	}
	assert.NotEmpty(s.T(), books, "Should have at least one book")

	testBook := books[0]

	// Step 2: Get chapters for the first book
	chapters, err := s.connector.GetBookChapters(testBook.URL, "")
	if err != nil {
		s.T().Skip("Skipping integration flow test: cannot get chapters")
		return
	}
	assert.NotEmpty(s.T(), chapters, "Should have at least one chapter")

	testChapter := chapters[0]

	// Step 3: Get content for the first chapter
	content, err := s.connector.GetChapterContent(testChapter.URL)
	if err != nil {
		s.T().Skip("Skipping integration flow test: cannot get content")
		return
	}

	assert.NoError(s.T(), err, "GetChapterContent should succeed")
	assert.Greater(s.T(), content.Count(), 0, "Chapter should have content")
}

// Benchmark tests for performance monitoring
func BenchmarkGetBooks(b *testing.B) {
	conn := connector.NewFreeWebNovelConnector()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := conn.GetBooks()
		if err != nil {
			b.Skipf("Benchmark skipped due to error: %v", err)
		}
	}
}

func BenchmarkGetBookCount(b *testing.B) {
	conn := connector.NewFreeWebNovelConnector()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := conn.GetBookCount()
		if err != nil {
			b.Skipf("Benchmark skipped due to error: %v", err)
		}
	}
}
