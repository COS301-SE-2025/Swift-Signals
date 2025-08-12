package intersection

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestDeleteIntersectionByID_Success() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock user intersection IDs stream
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStream.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStream.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.userClient.On("RemoveIntersectionID", ctx, userID, intersectionID).Return(nil, nil)
	suite.intrClient.On("DeleteIntersection", ctx, intersectionID).Return(nil, nil)

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().NoError(err)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_Forbidden() {
	userID := "valid-user-id"
	intersectionID := "intersection-999"

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock user intersection IDs stream
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStream.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStream.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("intersection not in user's intersection list", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_UserServiceError() {
	userID := "invalid-user-id"
	intersectionID := "intersection-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{}))

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_IntersectionServiceError() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock user intersection IDs stream
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStream.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStream.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.userClient.On("RemoveIntersectionID", ctx, userID, intersectionID).Return(nil, nil)
	suite.intrClient.On("DeleteIntersection", ctx, intersectionID).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_IntersectionNotFound() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock user intersection IDs stream
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStream.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStream.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.userClient.On("RemoveIntersectionID", ctx, userID, intersectionID).Return(nil, nil)
	suite.intrClient.On("DeleteIntersection", ctx, intersectionID).
		Return(nil, errs.NewNotFoundError("intersection not found", map[string]any{}))

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("intersection not found", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_UnassignError() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock user intersection IDs stream
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStream.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStream.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.userClient.On("RemoveIntersectionID", ctx, userID, intersectionID).
		Return(nil, errs.NewInternalError("failed to unassign intersection", nil, map[string]any{}))

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to unassign intersection", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_EmptyUserID() {
	userID := ""
	intersectionID := "intersection-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(nil, errs.NewValidationError("user ID cannot be empty", map[string]any{}))

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("user ID cannot be empty", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_EmptyIntersectionID() {
	userID := "valid-user-id"
	intersectionID := ""

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock user intersection IDs stream
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStream.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStream.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("intersection not in user's intersection list", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_ConflictError() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock user intersection IDs stream
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStream.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStream.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.userClient.On("RemoveIntersectionID", ctx, userID, intersectionID).Return(nil, nil)
	suite.intrClient.On("DeleteIntersection", ctx, intersectionID).
		Return(nil, errs.NewForbiddenError("intersection is currently being optimised and cannot be deleted", map[string]any{}))

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("intersection is currently being optimised and cannot be deleted", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersectionByID_StreamError() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock stream that returns an error
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	mockUserStream.On("Recv").
		Return(nil, errs.NewInternalError("stream error", nil, map[string]any{}))

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)

	err := suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("unable to retrieve intersection IDs", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func TestServiceDeleteIntersectionByID(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
