package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unluckythoughts/book-reader/server/models"
)

type FavoritesTestSuite struct {
	TestSuite
}

func TestFavoritesTestSuite(t *testing.T) {
	NewTestSuite(t)
}

func (suite *FavoritesTestSuite) getTestUserAndBook() (int, int, func()) {
	// Create a test user
	name := fmt.Sprintf("testuser_%d", time.Now().Unix())
	userRequest := &models.CreateUserRequest{
		Name: name,
	}
	user, err := suite.Client.CreateUser(userRequest)
	if err != nil {
		suite.T().Fatal("Failed to create test user:", err)
	}

	// Get a book
	booksResponse, err := suite.Client.ListBooks(1, 1, 0)
	if err != nil || len(booksResponse.Items) == 0 {
		suite.Client.DeleteUser(user.ID)
		suite.T().Fatal("No books available for testing")
	}

	bookID := booksResponse.Items[0].ID

	cleanup := func() {
		suite.Client.DeleteUser(user.ID)
	}

	return user.ID, bookID, cleanup
}

func (suite *FavoritesTestSuite) TestCreateFavorite() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Test creating a new favorite
	request := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}

	favorite, err := suite.Client.CreateFavorite(request)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), favorite)
	assert.NotZero(suite.T(), favorite.ID)
	assert.Equal(suite.T(), userID, favorite.UserID)
	assert.Equal(suite.T(), bookID, favorite.BookID)
	assert.False(suite.T(), favorite.UpdatedAt.IsZero())

	// Cleanup
	suite.Client.DeleteFavorite(favorite.ID)
}

func (suite *FavoritesTestSuite) TestListFavorites() {
	// Test listing favorites with pagination
	response, err := suite.Client.ListFavorites(1, 10, 0, 0)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))
	assert.LessOrEqual(suite.T(), len(response.Items), 10)
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 10, response.Pagination.Limit)
}

func (suite *FavoritesTestSuite) TestListFavoritesWithUserFilter() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Create a favorite
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}
	favorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteFavorite(favorite.ID)

	// Test listing favorites filtered by user
	response, err := suite.Client.ListFavorites(1, 10, userID, 0)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(1))

	// Verify all favorites belong to the specified user
	for _, fav := range response.Items {
		assert.Equal(suite.T(), userID, fav.UserID)
	}
}

func (suite *FavoritesTestSuite) TestListFavoritesWithBookFilter() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Create a favorite
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}
	favorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteFavorite(favorite.ID)

	// Test listing favorites filtered by book
	response, err := suite.Client.ListFavorites(1, 10, 0, bookID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(1))

	// Verify all favorites belong to the specified book
	for _, fav := range response.Items {
		assert.Equal(suite.T(), bookID, fav.BookID)
	}
}

func (suite *FavoritesTestSuite) TestListFavoritesWithBothFilters() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Create a favorite
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}
	favorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteFavorite(favorite.ID)

	// Test listing favorites filtered by both user and book
	response, err := suite.Client.ListFavorites(1, 10, userID, bookID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(1))

	// Verify all favorites match both filters
	for _, fav := range response.Items {
		assert.Equal(suite.T(), userID, fav.UserID)
		assert.Equal(suite.T(), bookID, fav.BookID)
	}
}

func (suite *FavoritesTestSuite) TestGetFavorite() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Create a favorite
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}
	createdFavorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteFavorite(createdFavorite.ID)

	// Test getting the favorite
	favorite, err := suite.Client.GetFavorite(createdFavorite.ID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), favorite)
	assert.Equal(suite.T(), createdFavorite.ID, favorite.ID)
	assert.Equal(suite.T(), userID, favorite.UserID)
	assert.Equal(suite.T(), bookID, favorite.BookID)
}

func (suite *FavoritesTestSuite) TestGetFavoriteNotFound() {
	// Test getting a non-existent favorite
	_, err := suite.Client.GetFavorite(999999)

	assert.Error(suite.T(), err)
}

func (suite *FavoritesTestSuite) TestUpdateFavorite() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Create a favorite
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}
	createdFavorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteFavorite(createdFavorite.ID)

	// Update the favorite
	updateRequest := &models.UpdateFavoriteRequest{
		Progress:   "5",
		Categories: "action,fantasy",
	}

	updatedFavorite, err := suite.Client.UpdateFavorite(createdFavorite.ID, updateRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedFavorite)
	assert.Equal(suite.T(), createdFavorite.ID, updatedFavorite.ID)
	assert.Equal(suite.T(), "5", updatedFavorite.Progress.String())
	assert.True(suite.T(), updatedFavorite.UpdatedAt.After(createdFavorite.UpdatedAt))
}

func (suite *FavoritesTestSuite) TestUpdateFavoriteProgress() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Create a favorite
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}
	createdFavorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteFavorite(createdFavorite.ID)

	// Update progress using PATCH
	progressRequest := &models.UpdateFavoriteProgressRequest{
		Chapter: 10,
		Level:   0,
	}

	updatedFavorite, err := suite.Client.UpdateFavoriteProgress(createdFavorite.ID, progressRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedFavorite)
	assert.Equal(suite.T(), createdFavorite.ID, updatedFavorite.ID)
}

func (suite *FavoritesTestSuite) TestDeleteFavorite() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// Create a favorite
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}
	createdFavorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)

	// Delete the favorite
	err = suite.Client.DeleteFavorite(createdFavorite.ID)
	assert.NoError(suite.T(), err)

	// Verify favorite is deleted
	_, err = suite.Client.GetFavorite(createdFavorite.ID)
	assert.Error(suite.T(), err)
}

func (suite *FavoritesTestSuite) TestDeleteFavoriteNotFound() {
	// Test deleting a non-existent favorite
	err := suite.Client.DeleteFavorite(999999)

	assert.Error(suite.T(), err)
}

func (suite *FavoritesTestSuite) TestFavoriteCRUDFlow() {
	userID, bookID, cleanup := suite.getTestUserAndBook()
	defer cleanup()

	// 1. Create
	createRequest := &models.CreateFavoriteRequest{
		UserID: userID,
		BookID: bookID,
	}

	favorite, err := suite.Client.CreateFavorite(createRequest)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), favorite)
	originalID := favorite.ID

	// 2. Read
	fetchedFavorite, err := suite.Client.GetFavorite(originalID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), originalID, fetchedFavorite.ID)

	// 3. Update
	updateRequest := &models.UpdateFavoriteRequest{
		Progress:   "15",
		Categories: "test",
	}

	updatedFavorite, err := suite.Client.UpdateFavorite(originalID, updateRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "15", updatedFavorite.Progress.String())

	// 4. Delete
	err = suite.Client.DeleteFavorite(originalID)
	assert.NoError(suite.T(), err)

	// 5. Verify deletion
	_, err = suite.Client.GetFavorite(originalID)
	assert.Error(suite.T(), err)
}

func (suite *FavoritesTestSuite) TestListFavoritesVerifyPreload() {
	// Verify that User and Book are preloaded in list response
	response, err := suite.Client.ListFavorites(1, 5, 0, 0)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)

	for _, favorite := range response.Items {
		assert.NotNil(suite.T(), favorite.User, "User should be preloaded for favorite ID %d", favorite.ID)
		assert.Equal(suite.T(), favorite.UserID, favorite.User.ID)

		assert.NotNil(suite.T(), favorite.Book, "Book should be preloaded for favorite ID %d", favorite.ID)
		assert.Equal(suite.T(), favorite.BookID, favorite.Book.ID)
	}
}
