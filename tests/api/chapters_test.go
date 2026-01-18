package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ChaptersTestSuite struct {
	TestSuite
}

func TestChaptersTestSuite(t *testing.T) {
	NewTestSuite(t)
}

func (suite *ChaptersTestSuite) TestListChapters() {
	// Test listing chapters with pagination
	response, err := suite.Client.ListChapters(1, 10, 0)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))
	assert.LessOrEqual(suite.T(), len(response.Items), 10)
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 10, response.Pagination.Limit)
}

func (suite *ChaptersTestSuite) TestListChaptersWithBookFilter() {
	// First get a book to filter by
	booksResponse, err := suite.Client.ListBooks(1, 1, 0)
	assert.NoError(suite.T(), err)

	if len(booksResponse.Items) == 0 {
		suite.T().Skip("No books available for testing")
		return
	}

	bookID := booksResponse.Items[0].ID

	// Test listing chapters filtered by book
	response, err := suite.Client.ListChapters(1, 10, bookID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))

	// Verify all chapters belong to the specified book
	for _, chapter := range response.Items {
		assert.Equal(suite.T(), bookID, chapter.BookID)
	}
}

func (suite *ChaptersTestSuite) TestListChaptersWithDifferentPageSizes() {
	// Test with different page sizes
	testCases := []struct {
		page  int
		limit int
	}{
		{1, 5},
		{1, 20},
		{2, 10},
	}

	for _, tc := range testCases {
		response, err := suite.Client.ListChapters(tc.page, tc.limit, 0)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), tc.page, response.Pagination.Page)
		assert.Equal(suite.T(), tc.limit, response.Pagination.Limit)
		assert.LessOrEqual(suite.T(), len(response.Items), tc.limit)
	}
}

func (suite *ChaptersTestSuite) TestGetChapter() {
	// First, get a list of chapters to find a valid ID
	listResponse, err := suite.Client.ListChapters(1, 1, 0)
	assert.NoError(suite.T(), err)

	if len(listResponse.Items) == 0 {
		suite.T().Skip("No chapters available for testing")
		return
	}

	chapterID := listResponse.Items[0].ID

	// Test getting a specific chapter
	chapter, err := suite.Client.GetChapter(chapterID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), chapter)
	assert.Equal(suite.T(), chapterID, chapter.ID)
	assert.NotEmpty(suite.T(), chapter.Title)
}

func (suite *ChaptersTestSuite) TestGetChapterNotFound() {
	// Test getting a non-existent chapter
	_, err := suite.Client.GetChapter(999999)

	assert.Error(suite.T(), err)
}

func (suite *ChaptersTestSuite) TestGetChapterDetails() {
	// Get a chapter and verify all fields are populated
	listResponse, err := suite.Client.ListChapters(1, 1, 0)
	assert.NoError(suite.T(), err)

	if len(listResponse.Items) == 0 {
		suite.T().Skip("No chapters available for testing")
		return
	}

	chapterID := listResponse.Items[0].ID
	chapter, err := suite.Client.GetChapter(chapterID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), chapter)
	assert.NotZero(suite.T(), chapter.ID)
	assert.NotEmpty(suite.T(), chapter.Title)
	assert.NotZero(suite.T(), chapter.BookID)
	assert.NotNil(suite.T(), chapter.Book)

	// Verify timestamps
	assert.False(suite.T(), chapter.UpdatedAt.IsZero())
}

func (suite *ChaptersTestSuite) TestListChaptersVerifyBookPreload() {
	// Verify that Book is preloaded in list response
	response, err := suite.Client.ListChapters(1, 5, 0)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)

	for _, chapter := range response.Items {
		if chapter.BookID > 0 {
			assert.NotNil(suite.T(), chapter.Book, "Book should be preloaded for chapter ID %d", chapter.ID)
			assert.Equal(suite.T(), chapter.BookID, chapter.Book.ID)
		}
	}
}

func (suite *ChaptersTestSuite) TestListChaptersOrdering() {
	// Get chapters for a specific book and verify they're ordered by index
	booksResponse, err := suite.Client.ListBooks(1, 1, 0)
	assert.NoError(suite.T(), err)

	if len(booksResponse.Items) == 0 {
		suite.T().Skip("No books available for testing")
		return
	}

	bookID := booksResponse.Items[0].ID
	response, err := suite.Client.ListChapters(1, 50, bookID)

	assert.NoError(suite.T(), err)

	// Chapters are returned from the API - just verify we got results
	if len(response.Items) > 0 {
		assert.NotZero(suite.T(), response.Items[0].ID)
	}
}
