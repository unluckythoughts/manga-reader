package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
)

type UsersTestSuite struct {
	TestSuite
}

func TestUsersTestSuite(t *testing.T) {
	NewTestSuite(t)
}

func (suite *UsersTestSuite) TestCreateUser() {
	// Test creating a new user
	name := fmt.Sprintf("testuser_%d", time.Now().Unix())
	request := &auth.RegisterRequest{
		Name:     name,
		Email:    fmt.Sprintf("%s@test.com", name),
		Password: "TestPass123!",
	}

	user, err := suite.Client.CreateUser(request)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotZero(suite.T(), user.ID)
	assert.Equal(suite.T(), request.Name, user.Name)
	assert.False(suite.T(), user.UpdatedAt.IsZero())

	// Cleanup
	suite.Client.DeleteUser(int(user.ID))
}

func (suite *UsersTestSuite) TestListUsers() {
	// Test listing users with pagination
	response, err := suite.Client.ListUsers(1, 10)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.GreaterOrEqual(suite.T(), response.Pagination.Total, int64(0))
	assert.LessOrEqual(suite.T(), len(response.Items), 10)
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 10, response.Pagination.Limit)
}

func (suite *UsersTestSuite) TestListUsersWithDifferentPageSizes() {
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
		response, err := suite.Client.ListUsers(tc.page, tc.limit)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), tc.page, response.Pagination.Page)
		assert.Equal(suite.T(), tc.limit, response.Pagination.Limit)
		assert.LessOrEqual(suite.T(), len(response.Items), tc.limit)
	}
}

func (suite *UsersTestSuite) TestGetUser() {
	// Create a user first
	name := fmt.Sprintf("testuser_%d", time.Now().Unix())
	createRequest := &auth.RegisterRequest{
		Name:     name,
		Email:    fmt.Sprintf("%s@test.com", name),
		Password: "TestPass123!",
	}

	createdUser, err := suite.Client.CreateUser(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteUser(int(createdUser.ID))

	// Test getting the user
	user, err := suite.Client.GetUser(int(createdUser.ID))

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), createdUser.ID, user.ID)
	assert.Equal(suite.T(), createdUser.Name, user.Name)
}

func (suite *UsersTestSuite) TestGetUserNotFound() {
	// Test getting a non-existent user
	_, err := suite.Client.GetUser(999999)

	assert.Error(suite.T(), err)
}

func (suite *UsersTestSuite) TestUpdateUser() {
	// Create a user first
	name := fmt.Sprintf("testuser_%d", time.Now().Unix())
	createRequest := &auth.RegisterRequest{
		Name:     name,
		Email:    fmt.Sprintf("%s@test.com", name),
		Password: "TestPass123!",
	}

	createdUser, err := suite.Client.CreateUser(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteUser(int(createdUser.ID))

	// Update the user
	newName := fmt.Sprintf("updated_%s", name)
	updateRequest := &auth.UpdateUserRequest{
		Name: newName,
	}

	updatedUser, err := suite.Client.UpdateUser(int(createdUser.ID), updateRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedUser)
	assert.Equal(suite.T(), createdUser.ID, updatedUser.ID)
	assert.Equal(suite.T(), newName, updatedUser.Name)
	assert.True(suite.T(), updatedUser.UpdatedAt.After(createdUser.UpdatedAt))
}

func (suite *UsersTestSuite) TestUpdateUserPartial() {
	// Create a user first
	name := fmt.Sprintf("testuser_%d", time.Now().Unix())
	createRequest := &auth.RegisterRequest{
		Name:     name,
		Email:    fmt.Sprintf("%s@test.com", name),
		Password: "TestPass123!",
	}

	createdUser, err := suite.Client.CreateUser(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteUser(int(createdUser.ID))

	// Update name
	newName := fmt.Sprintf("partial_%s", name)
	updateRequest := &auth.UpdateUserRequest{
		Name: newName,
	}

	updatedUser, err := suite.Client.UpdateUser(int(createdUser.ID), updateRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedUser)
	assert.Equal(suite.T(), newName, updatedUser.Name)
}

func (suite *UsersTestSuite) TestDeleteUser() {
	// Create a user first
	name := fmt.Sprintf("testuser_%d", time.Now().Unix())
	createRequest := &auth.RegisterRequest{
		Name:     name,
		Email:    fmt.Sprintf("%s@test.com", name),
		Password: "TestPass123!",
	}

	createdUser, err := suite.Client.CreateUser(createRequest)
	assert.NoError(suite.T(), err)

	// Delete the user
	err = suite.Client.DeleteUser(int(createdUser.ID))
	assert.NoError(suite.T(), err)

	// Verify user is deleted
	_, err = suite.Client.GetUser(int(createdUser.ID))
	assert.Error(suite.T(), err)
}

func (suite *UsersTestSuite) TestDeleteUserNotFound() {
	// Test deleting a non-existent user
	err := suite.Client.DeleteUser(999999)

	assert.Error(suite.T(), err)
}

func (suite *UsersTestSuite) TestCreateUserDuplicateEmail() {
	// Create a user
	name := fmt.Sprintf("testuser_%d", time.Now().Unix())
	email := fmt.Sprintf("%s@test.com", name)
	createRequest := &auth.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: "TestPass123!",
	}

	user1, err := suite.Client.CreateUser(createRequest)
	assert.NoError(suite.T(), err)
	defer suite.Client.DeleteUser(int(user1.ID))

	// Try to create another user with the same email
	createRequest2 := &auth.RegisterRequest{
		Name:     name + "_2",
		Email:    email, // Same email
		Password: "TestPass123!",
	}

	_, err = suite.Client.CreateUser(createRequest2)
	assert.Error(suite.T(), err) // Should fail due to duplicate email
}

func (suite *UsersTestSuite) TestUserCRUDFlow() {
	// Complete CRUD flow test

	// 1. Create
	name := fmt.Sprintf("crudtest_%d", time.Now().Unix())
	createRequest := &auth.RegisterRequest{
		Name:     name,
		Email:    fmt.Sprintf("%s@test.com", name),
		Password: "TestPass123!",
	}

	user, err := suite.Client.CreateUser(createRequest)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	originalID := int(user.ID)

	// 2. Read
	fetchedUser, err := suite.Client.GetUser(originalID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, fetchedUser.ID)
	assert.Equal(suite.T(), name, fetchedUser.Name)

	// 3. Update
	newName := fmt.Sprintf("updated_%s", name)
	updateRequest := &auth.UpdateUserRequest{
		Name: newName,
	}

	updatedUser, err := suite.Client.UpdateUser(originalID, updateRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, updatedUser.Name)

	// 4. Delete
	err = suite.Client.DeleteUser(originalID)
	assert.NoError(suite.T(), err)

	// 5. Verify deletion
	_, err = suite.Client.GetUser(originalID)
	assert.Error(suite.T(), err)
}
