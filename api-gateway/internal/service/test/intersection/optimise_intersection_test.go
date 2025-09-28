package intersection

import (
	"context"
	"io"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestOptimiseIntersectionByID_Success() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	expectedIntersection := createTestIntersection(
		intersectionID,
		"Test Intersection",
		"123 Test St",
		"Test City",
		"Test Province",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		0,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

	expectedOptimisationParams := &commonpb.OptimisationParameters{
		OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
		Parameters: &commonpb.SimulationParameters{
			IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

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
	suite.intrClient.On("GetIntersection", ctx, intersectionID).Return(expectedIntersection, nil)
	suite.optiClient.On("RunOptimisation", ctx, mock.AnythingOfType("model.OptimisationParameters")).
		Return(expectedOptimisationParams, nil)
	suite.intrClient.On("PutOptimisation", ctx, intersectionID, mock.AnythingOfType("model.OptimisationParameters")).
		Return(&intersectionpb.PutOptimisationResponse{Improved: true}, nil)

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().NoError(err)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	suite.optiClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_Forbidden() {
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

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("intersection not in user's intersection list", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_UserServiceError() {
	userID := "invalid-user-id"
	intersectionID := "intersection-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_OptimisationServiceError() {
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
	suite.intrClient.On("GetIntersection", ctx, intersectionID).
		Return(nil, errs.NewInternalError("optimisation service unavailable", nil, map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("optimisation service unavailable", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_IntersectionNotFound() {
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
	suite.intrClient.On("GetIntersection", ctx, intersectionID).
		Return(nil, errs.NewNotFoundError("intersection not found", map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("intersection not found", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_AlreadyOptimising() {
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

	expectedIntersection := createTestIntersection(
		intersectionID,
		"Test Intersection",
		"123 Test St",
		"Test City",
		"Test Province",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISING,
		0,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.intrClient.On("GetIntersection", ctx, intersectionID).Return(expectedIntersection, nil)
	suite.optiClient.On("RunOptimisation", ctx, mock.AnythingOfType("model.OptimisationParameters")).
		Return(nil, errs.NewConflictError("intersection is already being optimised", map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrConflict, svcError.Code)
	suite.Equal("intersection is already being optimised", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	suite.optiClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_InvalidParameters() {
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

	expectedIntersection := createTestIntersection(
		intersectionID,
		"Test Intersection",
		"123 Test St",
		"Test City",
		"Test Province",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		0,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.intrClient.On("GetIntersection", ctx, intersectionID).Return(expectedIntersection, nil)
	suite.optiClient.On("RunOptimisation", ctx, mock.AnythingOfType("model.OptimisationParameters")).
		Return(nil, errs.NewValidationError("invalid optimisation parameters", map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid optimisation parameters", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	suite.optiClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_EmptyUserID() {
	userID := ""
	intersectionID := "intersection-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(nil, errs.NewValidationError("user ID cannot be empty", map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("user ID cannot be empty", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_EmptyIntersectionID() {
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

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("intersection not in user's intersection list", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_StreamError() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Mock stream that returns an error
	mockUserStream := suite.NewMockUserIntersectionIDsStream()
	mockUserStream.On("Recv").
		Return(nil, errs.NewInternalError("stream error", nil, map[string]any{}))

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("unable to retrieve intersection IDs", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_UnauthorizedAccess() {
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

	expectedIntersection := createTestIntersection(
		intersectionID,
		"Test Intersection",
		"123 Test St",
		"Test City",
		"Test Province",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		0,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.intrClient.On("GetIntersection", ctx, intersectionID).Return(expectedIntersection, nil)
	suite.optiClient.On("RunOptimisation", ctx, mock.AnythingOfType("model.OptimisationParameters")).
		Return(nil, errs.NewUnauthorizedError("insufficient permissions for optimisation", map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("insufficient permissions for optimisation", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	suite.optiClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestOptimiseIntersectionByID_ServiceUnavailable() {
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

	expectedIntersection := createTestIntersection(
		intersectionID,
		"Test Intersection",
		"123 Test St",
		"Test City",
		"Test Province",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		0,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStream, nil)
	suite.intrClient.On("GetIntersection", ctx, intersectionID).Return(expectedIntersection, nil)
	suite.optiClient.On("RunOptimisation", ctx, mock.AnythingOfType("model.OptimisationParameters")).
		Return(nil, errs.NewUnavailableError("optimisation service is temporarily unavailable", map[string]any{}))

	err := suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)
	suite.Equal("optimisation service is temporarily unavailable", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	suite.optiClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}
