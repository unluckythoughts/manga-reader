package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unluckythoughts/book-reader/server/models"
)

type CategoriesTestSuite struct {
	TestSuite
}

func TestCategoriesTestSuite(t *testing.T) {
	NewTestSuite(t)
}

func (suite *CategoriesTestSuite) TestCreateCategory() {
	// Test creating a new category
	name := fmt.Sprintf("testcategory_%d", time.Now().Unix())
	request := &models.CreateCategoryRequest{
		Name: name,
	}

	category, err := suite.Client.CreateCategory(request)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), category)
	assert.NotZero(suite.T(), category.ID)
	assert.Equal(suite.T(), request.Name, category.Name)
	assert.False(suite.T(), category.UpdatedAt.IsZero())

	// Cleanup
	suite.Client.DeleteCategory(category.ID)
}

func (suite *CategoriesTestSuite) TestListCategories() {
	// Test listing categories with pagination
	response, err := suite.Client.ListCategories(1, 10)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))
	assert.LessOrEqual(suite.T(), len(response.Items), 10)
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 10, response.Pagination.Limit)
}

func (suite *CategoriesTestSuite) TestListCategoriesWithDifferentPageSizes() {
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
		response, err := suite.Client.ListCategories(tc.page, tc.limit)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), tc.page, response.Pagination.Page)
		assert.Equal(suite.T(), tc.limit, response.Pagination.Limit)
		assert.LessOrEqual(suite.T(), len(response.Items), tc.limit)
	}
}

func (suite *CategoriesTestSuite) TestGetCategory() {
	// Create a category first
	name := fmt.Sprintf("testcategory_%d", time.Now().Unix())
	createRequest := &models.CreateCategoryRequest{
		Name: name,
	}

	createdCategory, err := suite.Client.CreateCategory(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteCategory(createdCategory.ID)

	// Test getting the category
	category, err := suite.Client.GetCategory(createdCategory.ID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), category)
	assert.Equal(suite.T(), createdCategory.ID, category.ID)
	assert.Equal(suite.T(), createdCategory.Name, category.Name)
}

func (suite *CategoriesTestSuite) TestGetCategoryNotFound() {
	// Test getting a non-existent category
	_, err := suite.Client.GetCategory(999999)

	assert.Error(suite.T(), err)
}

func (suite *CategoriesTestSuite) TestUpdateCategory() {
	// Create a category first
	name := fmt.Sprintf("testcategory_%d", time.Now().Unix())
	createRequest := &models.CreateCategoryRequest{
		Name: name,
	}

	createdCategory, err := suite.Client.CreateCategory(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteCategory(createdCategory.ID)

	// Update the category
	newName := fmt.Sprintf("updated_%s", name)
	updateRequest := &models.UpdateCategoryRequest{
		Name: newName,
	}

	updatedCategory, err := suite.Client.UpdateCategory(createdCategory.ID, updateRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedCategory)
	assert.Equal(suite.T(), createdCategory.ID, updatedCategory.ID)
	assert.Equal(suite.T(), newName, updatedCategory.Name)
	assert.True(suite.T(), updatedCategory.UpdatedAt.After(createdCategory.UpdatedAt))
}

func (suite *CategoriesTestSuite) TestDeleteCategory() {
	// Create a category first
	name := fmt.Sprintf("testcategory_%d", time.Now().Unix())
	createRequest := &models.CreateCategoryRequest{
		Name: name,
	}

	createdCategory, err := suite.Client.CreateCategory(createRequest)
	assert.NoError(suite.T(), err)

	// Delete the category
	err = suite.Client.DeleteCategory(createdCategory.ID)
	assert.NoError(suite.T(), err)

	// Verify category is deleted
	_, err = suite.Client.GetCategory(createdCategory.ID)
	assert.Error(suite.T(), err)
}

func (suite *CategoriesTestSuite) TestDeleteCategoryNotFound() {
	// Test deleting a non-existent category
	err := suite.Client.DeleteCategory(999999)

	assert.Error(suite.T(), err)
}

func (suite *CategoriesTestSuite) TestCreateCategoryDuplicateName() {
	// Create a category
	name := fmt.Sprintf("testcategory_%d", time.Now().Unix())
	createRequest := &models.CreateCategoryRequest{
		Name: name,
	}

	category1, err := suite.Client.CreateCategory(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteCategory(category1.ID)

	// Try to create another category with the same name
	_, err = suite.Client.CreateCategory(createRequest)
	assert.Error(suite.T(), err) // Should fail due to duplicate name
}

func (suite *CategoriesTestSuite) TestCategoryCRUDFlow() {
	// Complete CRUD flow test

	// 1. Create
	name := fmt.Sprintf("crudtest_%d", time.Now().Unix())
	createRequest := &models.CreateCategoryRequest{
		Name: name,
	}

	category, err := suite.Client.CreateCategory(createRequest)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), category)
	originalID := category.ID

	// 2. Read
	fetchedCategory, err := suite.Client.GetCategory(originalID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), originalID, fetchedCategory.ID)
	assert.Equal(suite.T(), name, fetchedCategory.Name)

	// 3. Update
	newName := fmt.Sprintf("updated_%s", name)
	updateRequest := &models.UpdateCategoryRequest{
		Name: newName,
	}

	updatedCategory, err := suite.Client.UpdateCategory(originalID, updateRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, updatedCategory.Name)

	// 4. Delete
	err = suite.Client.DeleteCategory(originalID)
	assert.NoError(suite.T(), err)

	// 5. Verify deletion
	_, err = suite.Client.GetCategory(originalID)
	assert.Error(suite.T(), err)
}

func (suite *CategoriesTestSuite) TestListCategoriesVerifyFields() {
	// Verify that all required fields are present in list response
	response, err := suite.Client.ListCategories(1, 5)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)

	for _, category := range response.Items {
		assert.NotZero(suite.T(), category.ID, "Category ID should not be zero")
		assert.NotEmpty(suite.T(), category.Name, "Category name should not be empty")
		assert.False(suite.T(), category.UpdatedAt.IsZero(), "UpdatedAt should be set")
	}
}
