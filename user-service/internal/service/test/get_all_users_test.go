package test

import (
	"context"
	"errors"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestGetAllUsers_Success() {
	page := int32(1)
	pageSize := int32(10)
	filter := ""
	expectedUsers := []*model.User{
		{
			ID:      "550e8400-e29b-41d4-a716-446655440000",
			Name:    "User One",
			Email:   "user1@example.com",
			IsAdmin: false,
		},
		{
			ID:      "550e8400-e29b-41d4-a716-446655440001",
			Name:    "User Two",
			Email:   "user2@example.com",
			IsAdmin: true,
		},
	}

	suite.repo.On("ListUsers", mock.Anything, 10, 0).Return(expectedUsers, nil)

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Len(result, 2)
	suite.Equal(expectedUsers, result)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_WithPagination() {
	page := int32(2)
	pageSize := int32(5)
	filter := ""
	expectedUsers := []*model.User{
		{
			ID:      "550e8400-e29b-41d4-a716-446655440005",
			Name:    "User Six",
			Email:   "user6@example.com",
			IsAdmin: false,
		},
	}

	suite.repo.On("ListUsers", mock.Anything, 5, 5).Return(expectedUsers, nil)

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Len(result, 1)
	suite.Equal(expectedUsers, result)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_EmptyResults() {
	page := int32(1)
	pageSize := int32(10)
	filter := ""
	expectedUsers := []*model.User{}

	suite.repo.On("ListUsers", mock.Anything, 10, 0).Return(expectedUsers, nil)

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Empty(result)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_InvalidPaginationNegativePage() {
	page := int32(-1)
	pageSize := int32(10)
	filter := ""

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"page": "Page must be at least 1 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_InvalidPaginationZeroPage() {
	page := int32(0)
	pageSize := int32(10)
	filter := ""

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"page": "Page must be at least 1 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_InvalidPaginationNegativePageSize() {
	page := int32(1)
	pageSize := int32(-1)
	filter := ""

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"pagesize": "PageSize must be at least 1 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_InvalidPaginationZeroPageSize() {
	page := int32(1)
	pageSize := int32(0)
	filter := ""

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"pagesize": "PageSize must be at least 1 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_InvalidPaginationExcessivePageSize() {
	page := int32(1)
	pageSize := int32(200)
	filter := ""

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"pagesize": "PageSize must be at most 100 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_FilterTooLong() {
	page := int32(1)
	pageSize := int32(10)
	filter := "a very long filter string that exceeds the maximum allowed length of 255 characters and should trigger a validation error when the method is properly implemented with validation logic that checks the filter parameter length against the defined constraints in the contract struct validation tags"

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"filter": "Filter must be at most 255 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_RepositoryError() {
	page := int32(1)
	pageSize := int32(10)
	filter := ""
	repoError := errors.New("database connection failed")

	suite.repo.On("ListUsers", mock.Anything, 10, 0).Return(nil, repoError)

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to retrieve users", svcError.Message)
	expectedContext := map[string]any{
		"page":     page,
		"pageSize": pageSize,
		"filter":   "",
	}
	suite.Equal(expectedContext, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_ValidMinimumParameters() {
	page := int32(1)
	pageSize := int32(1)
	filter := ""
	expectedUsers := []*model.User{
		{
			ID:      "550e8400-e29b-41d4-a716-446655440000",
			Name:    "Single User",
			Email:   "single@example.com",
			IsAdmin: false,
		},
	}

	suite.repo.On("ListUsers", mock.Anything, 1, 0).Return(expectedUsers, nil)

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Len(result, 1)
	suite.Equal(expectedUsers, result)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_ValidMaximumPageSize() {
	page := int32(1)
	pageSize := int32(100)
	filter := ""
	expectedUsers := []*model.User{}

	suite.repo.On("ListUsers", mock.Anything, 100, 0).Return(expectedUsers, nil)

	ctx := context.Background()

	result, err := suite.service.GetAllUsers(ctx, page, pageSize, filter)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Empty(result)

	suite.repo.AssertExpectations(suite.T())
}
