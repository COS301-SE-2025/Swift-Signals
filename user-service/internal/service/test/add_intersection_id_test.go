package test

import (
	"context"
	"errors"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestAddIntersectionID_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	existingIntersectionIDs := []string{"intersection-2", "intersection-3"}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(existingIntersectionIDs, nil)
	suite.repo.On("AddIntersectionID", mock.Anything, userID, intersectionID).Return(nil)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_FirstIntersection() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	existingIntersectionIDs := []string{}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(existingIntersectionIDs, nil)
	suite.repo.On("AddIntersectionID", mock.Anything, userID, intersectionID).Return(nil)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_GetUserByIDError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_GetUserByIDServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_IntersectionAlreadyExists() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	existingIntersectionIDs := []string{
		"intersection-1",
		"intersection-2",
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(existingIntersectionIDs, nil)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_GetIntersectionsError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	repoError := errors.New("failed to fetch intersections")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to check existing intersection IDs", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_GetIntersectionsServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	serviceError := errs.NewInternalError(
		"failed to fetch intersections",
		errors.New("db error"),
		nil,
	)

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to fetch intersections", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_AddIntersectionError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	existingIntersectionIDs := []string{"intersection-2", "intersection-3"}
	addError := errors.New("failed to add intersection")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(existingIntersectionIDs, nil)
	suite.repo.On("AddIntersectionID", mock.Anything, userID, intersectionID).Return(addError)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to add intersection ID", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_AddIntersectionServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := "intersection-1"
	expectedUser := &model.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		IsAdmin: false,
	}
	existingIntersectionIDs := []string{"intersection-2", "intersection-3"}
	serviceError := errs.NewInternalError("failed to add intersection", errors.New("db error"), nil)

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(existingIntersectionIDs, nil)
	suite.repo.On("AddIntersectionID", mock.Anything, userID, intersectionID).Return(serviceError)

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to add intersection", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_EmptyUserID() {
	userID := ""
	intersectionID := "intersection-1"

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_EmptyIntersectionID() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionID := ""

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestAddIntersectionID_InvalidUUIDFormat() {
	userID := "invalid-uuid-format"
	intersectionID := "intersection-1"

	ctx := context.Background()

	err := suite.service.AddIntersectionID(ctx, userID, intersectionID)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func TestServiceAddIntersectionID(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
