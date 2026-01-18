package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type BooksTestSuite struct {
	TestSuite
}

func TestBooksTestSuite(t *testing.T) {
	NewTestSuite(t)
}

func (suite *BooksTestSuite) TestListBooks() {
	// Test listing books with pagination
	response, err := suite.Client.ListBooks(1, 10, 0)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))
	assert.LessOrEqual(suite.T(), len(response.Items), 10)
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 10, response.Pagination.Limit)
}

func (suite *BooksTestSuite) TestListBooksWithSourceFilter() {
	// Test listing books filtered by source
	response, err := suite.Client.ListBooks(1, 10, 1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))

	// Verify all books belong to the specified source
	for _, book := range response.Items {
		assert.Equal(suite.T(), 1, book.SourceID)
	}
}

func (suite *BooksTestSuite) TestListBooksWithDifferentPageSizes() {
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
		response, err := suite.Client.ListBooks(tc.page, tc.limit, 0)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), tc.page, response.Pagination.Page)
		assert.Equal(suite.T(), tc.limit, response.Pagination.Limit)
		assert.LessOrEqual(suite.T(), len(response.Items), tc.limit)
	}
}

func (suite *BooksTestSuite) TestGetBook() {
	// First, get a list of books to find a valid ID
	listResponse, err := suite.Client.ListBooks(1, 1, 0)
	assert.NoError(suite.T(), err)

	if len(listResponse.Items) == 0 {
		suite.T().Skip("No books available for testing")
		return
	}

	bookID := listResponse.Items[0].ID

	// Test getting a specific book
	book, err := suite.Client.GetBook(bookID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), book)
	assert.Equal(suite.T(), bookID, book.ID)
	assert.NotEmpty(suite.T(), book.Title)
}

func (suite *BooksTestSuite) TestGetBookNotFound() {
	// Test getting a non-existent book
	_, err := suite.Client.GetBook(999999)

	assert.Error(suite.T(), err)
}

func (suite *BooksTestSuite) TestGetBookDetails() {
	// Get a book and verify all fields are populated
	listResponse, err := suite.Client.ListBooks(1, 1, 0)
	assert.NoError(suite.T(), err)

	if len(listResponse.Items) == 0 {
		suite.T().Skip("No books available for testing")
		return
	}

	bookID := listResponse.Items[0].ID
	book, err := suite.Client.GetBook(bookID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), book)
	assert.NotZero(suite.T(), book.ID)
	assert.NotEmpty(suite.T(), book.Title)
	assert.NotZero(suite.T(), book.SourceID)
	assert.NotNil(suite.T(), book.Source)

	// Verify timestamps
	assert.False(suite.T(), book.UpdatedAt.IsZero())
}

func (suite *BooksTestSuite) TestListBooksVerifySourcePreload() {
	// Verify that Source is preloaded in list response
	response, err := suite.Client.ListBooks(1, 5, 0)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)

	for _, book := range response.Items {
		if book.SourceID > 0 {
			assert.NotNil(suite.T(), book.Source, "Source should be preloaded for book ID %d", book.ID)
			assert.Equal(suite.T(), book.SourceID, book.Source.ID)
		}
	}
}
