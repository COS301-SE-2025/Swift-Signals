package intersection

import (
	"context"
	"io"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
)

func (suite *TestSuite) TestUpdateIntersectionByID_Success() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	request := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

	expectedIntersectionIDs := []string{"intersection-123", "intersection-456"}

	expectedUpdatedIntersection := createTestIntersection(
		"intersection-123",
		"Updated Intersection",
		"456 Updated Street",
		"Johannesburg",
		"Gauteng",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISED,
		5,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

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
	suite.intrClient.On("UpdateIntersection", ctx, intersectionID, "Updated Intersection", model.Details{
		Address:  "456 Updated Street",
		City:     "Johannesburg",
		Province: "Gauteng",
	}).
		Return(expectedUpdatedIntersection, nil)

	result, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().NoError(err)
	suite.Equal("intersection-123", result.ID)
	suite.Equal("Updated Intersection", result.Name)
	suite.Equal("456 Updated Street", result.Details.Address)
	suite.Equal("Johannesburg", result.Details.City)
	suite.Equal("Gauteng", result.Details.Province)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersectionByID_Forbidden() {
	userID := "valid-user-id"
	intersectionID := "intersection-999"

	request := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

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

	_, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("intersection not in user's intersection list", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersectionByID_UserServiceError() {
	userID := "invalid-user-id"
	intersectionID := "intersection-123"

	request := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(nil, errs.NewNotFoundError("user not found", map[string]any{}))

	_, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersectionByID_IntersectionServiceError() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	request := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

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
	suite.intrClient.On("UpdateIntersection", ctx, intersectionID, "Updated Intersection", model.Details{
		Address:  "456 Updated Street",
		City:     "Johannesburg",
		Province: "Gauteng",
	}).
		Return(nil, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	_, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("database connection failed", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersectionByID_IntersectionNotFound() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	request := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

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
	suite.intrClient.On("UpdateIntersection", ctx, intersectionID, "Updated Intersection", model.Details{
		Address:  "456 Updated Street",
		City:     "Johannesburg",
		Province: "Gauteng",
	}).
		Return(nil, errs.NewNotFoundError("intersection not found", map[string]any{}))

	_, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("intersection not found", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersectionByID_ValidationError() {
	userID := "valid-user-id"
	intersectionID := "intersection-123"

	request := model.UpdateIntersectionRequest{
		Name: "", // Empty name should cause validation error
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

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
	suite.intrClient.On("UpdateIntersection", ctx, intersectionID, "", model.Details{
		Address:  "456 Updated Street",
		City:     "Johannesburg",
		Province: "Gauteng",
	}).Return(nil, errs.NewValidationError("intersection name cannot be empty", map[string]any{}))

	_, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("intersection name cannot be empty", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	suite.intrClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersectionByID_EmptyUserID() {
	userID := ""
	intersectionID := "intersection-123"

	request := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(nil, errs.NewValidationError("user ID cannot be empty", map[string]any{}))

	_, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("user ID cannot be empty", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersectionByID_EmptyIntersectionID() {
	userID := "valid-user-id"
	intersectionID := ""

	request := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "456 Updated Street",
			City:     "Johannesburg",
			Province: "Gauteng",
		},
	}

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

	_, err := suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, request)

	suite.Require().Error(err)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)
	suite.Equal("intersection not in user's intersection list", svcError.Message)

	suite.userClient.AssertExpectations(suite.T())
	mockUserStream.AssertExpectations(suite.T())
}
