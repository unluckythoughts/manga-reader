package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SourcesTestSuite struct {
	TestSuite
}

func TestSourcesTestSuite(t *testing.T) {
	NewTestSuite(t)
}

func (suite *SourcesTestSuite) TestListSources() {
	// Test listing sources with pagination
	response, err := suite.Client.ListSources(1, 10)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))
	assert.LessOrEqual(suite.T(), len(response.Items), 10)
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 10, response.Pagination.Limit)
}

func (suite *SourcesTestSuite) TestListSourcesWithDifferentPageSizes() {
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
		response, err := suite.Client.ListSources(tc.page, tc.limit)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), tc.page, response.Pagination.Page)
		assert.Equal(suite.T(), tc.limit, response.Pagination.Limit)
		assert.LessOrEqual(suite.T(), len(response.Items), tc.limit)
	}
}

func (suite *SourcesTestSuite) TestGetSource() {
	// First, get a list of sources to find a valid ID
	listResponse, err := suite.Client.ListSources(1, 1)
	assert.NoError(suite.T(), err)

	if len(listResponse.Items) == 0 {
		suite.T().Skip("No sources available for testing")
		return
	}

	sourceID := listResponse.Items[0].ID

	// Test getting a specific source
	source, err := suite.Client.GetSource(sourceID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), source)
	assert.Equal(suite.T(), sourceID, source.ID)
	assert.NotEmpty(suite.T(), source.Name)
}

func (suite *SourcesTestSuite) TestGetSourceNotFound() {
	// Test getting a non-existent source
	_, err := suite.Client.GetSource(999999)

	assert.Error(suite.T(), err)
}

func (suite *SourcesTestSuite) TestGetSourceDetails() {
	// Get a source and verify all fields are populated
	listResponse, err := suite.Client.ListSources(1, 1)
	assert.NoError(suite.T(), err)

	if len(listResponse.Items) == 0 {
		suite.T().Skip("No sources available for testing")
		return
	}

	sourceID := listResponse.Items[0].ID
	source, err := suite.Client.GetSource(sourceID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), source)
	assert.NotZero(suite.T(), source.ID)
	assert.NotEmpty(suite.T(), source.Name)
	assert.NotEmpty(suite.T(), source.Domain)

	// Verify timestamps
	assert.False(suite.T(), source.UpdatedAt.IsZero())
}

func (suite *SourcesTestSuite) TestListSourcesVerifyFields() {
	// Verify that all required fields are present in list response
	response, err := suite.Client.ListSources(1, 5)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)

	for _, source := range response.Items {
		assert.NotZero(suite.T(), source.ID, "Source ID should not be zero")
		assert.NotEmpty(suite.T(), source.Name, "Source name should not be empty")
		assert.NotEmpty(suite.T(), source.Domain, "Source domain should not be empty")
		assert.False(suite.T(), source.UpdatedAt.IsZero(), "UpdatedAt should be set")
	}
}

func (suite *SourcesTestSuite) TestListAllSources() {
	// Test fetching all sources by using a large limit
	response, err := suite.Client.ListSources(1, 100)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), int64(len(response.Items)), response.Pagination.Total)
}
