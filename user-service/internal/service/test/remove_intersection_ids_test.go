package test

import (
	"context"
	"errors"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestRemoveIntersectionIDs_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1", "intersection-2"}

	existingUser := &model.User{
		ID:              userID,
		Name:            "Test User",
		Email:           "test@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"intersection-1", "intersection-2", "intersection-3"},
	}

	currentIntersectionIDs := []string{"intersection-1", "intersection-2", "intersection-3"}

	updatedUser := &model.User{
		ID:              userID,
		Name:            "Test User",
		Email:           "test@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"intersection-3"},
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(currentIntersectionIDs, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && len(u.IntersectionIDs) == 1 &&
			u.IntersectionIDs[0] == "intersection-3"
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_SingleIntersection() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}

	existingUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-1", "intersection-2"},
	}

	currentIntersectionIDs := []string{"intersection-1", "intersection-2"}

	updatedUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-2"},
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(currentIntersectionIDs, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && len(u.IntersectionIDs) == 1 &&
			u.IntersectionIDs[0] == "intersection-2"
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_RemoveAllIntersections() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1", "intersection-2"}

	existingUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-1", "intersection-2"},
	}

	currentIntersectionIDs := []string{"intersection-1", "intersection-2"}

	updatedUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{},
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(currentIntersectionIDs, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && len(u.IntersectionIDs) == 0
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_InvalidUserID() {
	userID := "invalid-uuid"
	intersectionIDs := []string{"intersection-1"}

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID must be a valid UUID",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_EmptyUserID() {
	userID := ""
	intersectionIDs := []string{"intersection-1"}

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_EmptyIntersectionIDsList() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{}

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"intersectionids": "IntersectionIDs must be at least 1 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_EmptyIntersectionID() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"", "intersection-1"}

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_UserNotFound() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}
	userNotFoundError := errs.NewNotFoundError("user not found", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, userNotFoundError)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_NonExistentIntersection() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"nonexistent-intersection", "intersection-1"}

	existingUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-1", "intersection-2"},
	}

	currentIntersectionIDs := []string{"intersection-1", "intersection-2"}

	updatedUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-2"},
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(currentIntersectionIDs, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && len(u.IntersectionIDs) == 1 &&
			u.IntersectionIDs[0] == "intersection-2"
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_RepositoryGetUserError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_RepositoryGetIntersectionsError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}

	existingUser := &model.User{
		ID: userID,
	}
	repoError := errors.New("failed to fetch intersections")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to fetch current intersections", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_RepositoryGetIntersectionsServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}

	existingUser := &model.User{
		ID: userID,
	}
	repoError := errs.NewUnauthorizedError("failed to remove intersections", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("failed to remove intersections", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_RepositoryUpdateError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}

	existingUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-1", "intersection-2"},
	}

	currentIntersectionIDs := []string{"intersection-1", "intersection-2"}
	repoError := errors.New("database update failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(currentIntersectionIDs, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_RepositoryUpdateServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}

	existingUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-1", "intersection-2"},
	}

	currentIntersectionIDs := []string{"intersection-1", "intersection-2"}
	repoError := errs.NewUnauthorizedError("failed to update user", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(currentIntersectionIDs, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, repoError)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_ServiceErrorPropagation() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1"}
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRemoveIntersectionIDs_DuplicateIntersectionIDs() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	intersectionIDs := []string{"intersection-1", "intersection-1", "intersection-2"}

	existingUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-1", "intersection-2", "intersection-3"},
	}

	currentIntersectionIDs := []string{"intersection-1", "intersection-2", "intersection-3"}

	updatedUser := &model.User{
		ID:              userID,
		IntersectionIDs: []string{"intersection-3"},
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetIntersectionsByUserID", mock.Anything, userID).
		Return(currentIntersectionIDs, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && len(u.IntersectionIDs) == 1 &&
			u.IntersectionIDs[0] == "intersection-3"
	})).Return(updatedUser, nil)

	ctx := context.Background()

	err := suite.service.RemoveIntersectionIDs(ctx, userID, intersectionIDs)

	suite.Require().NoError(err)

	suite.repo.AssertExpectations(suite.T())
}
