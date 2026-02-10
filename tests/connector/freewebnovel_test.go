package connector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unluckythoughts/book-reader/server/connector"
	"github.com/unluckythoughts/book-reader/server/models"
)

type FreeWebNovelTestSuite struct {
	ConnectorTestSuite
	connector models.IConnector
}

func TestFreeWebNovelTestSuite(t *testing.T) {
	suite := &FreeWebNovelTestSuite{}
	suite.SetupSuite()
	suite.RunTests(t)
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

func (s *FreeWebNovelTestSuite) RunTests(t *testing.T) {
	t.Run("TestConnectorProperties", s.TestConnectorProperties)
	t.Run("TestGetName", s.TestGetName)
	t.Run("TestGetDomain", s.TestGetDomain)
	t.Run("TestGetIconURL", s.TestGetIconURL)
	t.Run("TestGetSelectors", s.TestGetSelectors)
	t.Run("TestGetBooks", s.TestGetBooks)
	t.Run("TestGetBookCount", s.TestGetBookCount)
	t.Run("TestGetBookChapters", s.TestGetBookChapters)
	t.Run("TestGetBookChaptersWithFilter", s.TestGetBookChaptersWithFilter)
	t.Run("TestGetChapterContent", s.TestGetChapterContent)
	t.Run("TestIntegrationFlow", s.TestIntegrationFlow)
}

func (s *FreeWebNovelTestSuite) TestConnectorProperties(t *testing.T) {
	assert.NotNil(t, s.connector, "Connector should not be nil")
}

func (s *FreeWebNovelTestSuite) TestGetName(t *testing.T) {
	name := s.connector.GetName()
	assert.Equal(t, "FreeWebNovel", name, "Connector name should be 'FreeWebNovel'")
}

func (s *FreeWebNovelTestSuite) TestGetDomain(t *testing.T) {
	domain := s.connector.GetDomain()
	assert.Equal(t, "freewebnovel.com", domain, "Domain should be 'freewebnovel.com'")
}

func (s *FreeWebNovelTestSuite) TestGetIconURL(t *testing.T) {
	iconURL := s.connector.GetIconURL()
	assert.Equal(t, "/static/freewebnovel/favicon.ico", iconURL, "Icon URL should match expected path")
}

func (s *FreeWebNovelTestSuite) TestGetSelectors(t *testing.T) {
	selectors := s.connector.GetSelectors()

	assert.NotEmpty(t, selectors.BookListItem, "BookListItem selector should not be empty")
	assert.NotEmpty(t, selectors.LastPage, "LastPage selector should not be empty")
	assert.NotEmpty(t, selectors.NextPageURLPattern, "NextPageURLPattern should not be empty")

	// Test book selectors
	assert.NotEmpty(t, selectors.Book.URL, "Book URL selector should not be empty")
	assert.NotEmpty(t, selectors.Book.Title, "Book Title selector should not be empty")
	assert.NotEmpty(t, selectors.Book.CoverImage, "Book CoverImage selector should not be empty")
	assert.NotEmpty(t, selectors.Book.Synopsis, "Book Synopsis selector should not be empty")
	assert.NotEmpty(t, selectors.Book.ChapterListItem, "ChapterListItem selector should not be empty")

	// Test chapter selectors
	assert.NotEmpty(t, selectors.Book.Chapter.URL, "Chapter URL selector should not be empty")
	assert.NotEmpty(t, selectors.Book.Chapter.Number, "Chapter Number selector should not be empty")
	assert.NotEmpty(t, selectors.Book.Chapter.Title, "Chapter Title selector should not be empty")
	assert.NotEmpty(t, selectors.Book.Chapter.Content.Data, "Chapter Content Data selector should not be empty")
}

func (s *FreeWebNovelTestSuite) TestGetBooks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	books, err := s.connector.GetBooks()

	if err != nil {
		t.Logf("Warning: GetBooks failed (site may be down or blocking requests): %v", err)
		if s.isRodBlockedError(err) {
			t.Skip("Skipping test: Windows is blocking browser automation (leakless.exe). See README for instructions.")
		}
		t.Skip("Skipping test due to network/site availability issues")
		return
	}

	require.NoError(t, err, "GetBooks should not return an error")
	assert.NotEmpty(t, books, "Books list should not be empty")

	// Validate first book structure
	if len(books) > 0 {
		book := books[0]
		assert.NotEmpty(t, book.Title, "Book title should not be empty")
		assert.NotEmpty(t, book.URL, "Book URL should not be empty")
		assert.Contains(t, book.URL, "/", "Book URL should be a valid path")

		t.Logf("Sample book: Title='%s', URL='%s'", book.Title, book.URL)
	}
}

func (s *FreeWebNovelTestSuite) TestGetBookCount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	count, books, err := s.connector.GetBookCount()

	if err != nil {
		t.Logf("Warning: GetBookCount failed (site may be down or blocking requests): %v", err)
		if s.isRodBlockedError(err) {
			t.Skip("Skipping test: Windows is blocking browser automation (leakless.exe). See README for instructions.")
		}
		t.Skip("Skipping test due to network/site availability issues")
		return
	}

	require.NoError(t, err, "GetBookCount should not return an error")
	assert.Greater(t, count, 0, "Book count should be greater than 0")

	// books may be nil if last page selector is present
	if books != nil {
		assert.NotEmpty(t, books, "Books list should not be empty when returned")
	}

	t.Logf("Total book count: %d", count)
}

func (s *FreeWebNovelTestSuite) TestGetBookChapters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// First get a book to test with
	books, err := s.connector.GetBooks()
	if err != nil {
		t.Skip("Skipping test: cannot get books for testing")
		return
	}

	if len(books) == 0 {
		t.Skip("Skipping test: no books available")
		return
	}

	bookURL := books[0].URL
	t.Logf("Testing with book URL: %s", bookURL)

	chapters, err := s.connector.GetBookChapters(bookURL, "")

	if err != nil {
		t.Logf("Warning: GetBookChapters failed: %v", err)
		t.Skip("Skipping test due to network/site availability issues")
		return
	}

	require.NoError(t, err, "GetBookChapters should not return an error")
	assert.NotEmpty(t, chapters, "Chapters list should not be empty")

	// Validate first chapter structure
	if len(chapters) > 0 {
		chapter := chapters[0]
		assert.NotEmpty(t, chapter.Title, "Chapter title should not be empty")
		assert.NotEmpty(t, chapter.URL, "Chapter URL should not be empty")
		assert.NotEmpty(t, chapter.Number, "Chapter number should not be empty")

		t.Logf("Sample chapter: Number='%s', Title='%s', URL='%s'",
			chapter.Number, chapter.Title, chapter.URL)
	}
}

func (s *FreeWebNovelTestSuite) TestGetBookChaptersWithFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// First get a book to test with
	books, err := s.connector.GetBooks()
	if err != nil {
		t.Skip("Skipping test: cannot get books for testing")
		return
	}

	if len(books) == 0 {
		t.Skip("Skipping test: no books available")
		return
	}

	bookURL := books[0].URL

	// Get all chapters first
	allChapters, err := s.connector.GetBookChapters(bookURL, "")
	if err != nil || len(allChapters) == 0 {
		t.Skip("Skipping test: cannot get chapters for testing")
		return
	}

	// Try to filter chapters greater than chapter 1
	filteredChapters, err := s.connector.GetBookChapters(bookURL, "1")

	if err != nil {
		t.Logf("Warning: GetBookChapters with filter failed: %v", err)
		t.Skip("Skipping test due to network/site availability issues")
		return
	}

	require.NoError(t, err, "GetBookChapters with filter should not return an error")

	// Filtered chapters should be less than or equal to all chapters
	assert.LessOrEqual(t, len(filteredChapters), len(allChapters),
		"Filtered chapters should be less than or equal to all chapters")

	t.Logf("All chapters: %d, Filtered (>1): %d", len(allChapters), len(filteredChapters))
}

func (s *FreeWebNovelTestSuite) TestGetChapterContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// First get a book and then a chapter
	books, err := s.connector.GetBooks()
	if err != nil {
		t.Skip("Skipping test: cannot get books for testing")
		return
	}

	if len(books) == 0 {
		t.Skip("Skipping test: no books available")
		return
	}

	chapters, err := s.connector.GetBookChapters(books[0].URL, "")
	if err != nil || len(chapters) == 0 {
		t.Skip("Skipping test: cannot get chapters for testing")
		return
	}

	chapterURL := chapters[0].URL
	t.Logf("Testing chapter content with URL: %s", chapterURL)

	content, err := s.connector.GetChapterContent(chapterURL)

	if err != nil {
		t.Logf("Warning: GetChapterContent failed: %v", err)
		t.Skip("Skipping test due to network/site availability issues")
		return
	}

	require.NoError(t, err, "GetChapterContent should not return an error")
	assert.NotEmpty(t, content.Values(), "Chapter content should not be empty")
	assert.Greater(t, content.Count(), 0, "Content should have at least one paragraph")

	// Validate content structure
	contentValues := content.Values()
	for i, paragraph := range contentValues {
		if i < 3 { // Log first 3 paragraphs
			t.Logf("Paragraph %d length: %d chars", i+1, len(paragraph))
		}
		assert.NotEmpty(t, paragraph, "Each content paragraph should not be empty")
	}

	t.Logf("Total content paragraphs: %d", content.Count())
}

func (s *FreeWebNovelTestSuite) TestIntegrationFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test the complete flow: Books -> Chapters -> Content

	// Step 1: Get books
	books, err := s.connector.GetBooks()
	if err != nil {
		t.Skip("Skipping integration flow test: cannot get books")
		return
	}
	require.NotEmpty(t, books, "Should have at least one book")

	testBook := books[0]
	t.Logf("Testing integration flow with book: '%s'", testBook.Title)

	// Step 2: Get chapters for the first book
	chapters, err := s.connector.GetBookChapters(testBook.URL, "")
	if err != nil {
		t.Skip("Skipping integration flow test: cannot get chapters")
		return
	}
	require.NotEmpty(t, chapters, "Should have at least one chapter")

	testChapter := chapters[0]
	t.Logf("Testing with chapter: Number='%s', Title='%s'",
		testChapter.Number, testChapter.Title)

	// Step 3: Get content for the first chapter
	content, err := s.connector.GetChapterContent(testChapter.URL)
	if err != nil {
		t.Skip("Skipping integration flow test: cannot get content")
		return
	}

	require.NoError(t, err, "GetChapterContent should succeed")
	assert.Greater(t, content.Count(), 0, "Chapter should have content")

	t.Logf("✓ Integration flow test completed successfully")
	t.Logf("  Book: %s", testBook.Title)
	t.Logf("  Chapter: %s - %s", testChapter.Number, testChapter.Title)
	t.Logf("  Content paragraphs: %d", content.Count())
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
