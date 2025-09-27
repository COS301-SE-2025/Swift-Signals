package intersection

import (
	"context"
	"io"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
	optimisationpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/optimisation/v1"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

// TestIntegrationSuite tests the intersection service with integrated workflows
type TestIntegrationSuite struct {
	TestSuite
}

func (suite *TestIntegrationSuite) TestCompleteIntersectionLifecycle() {
	userID := "integration-user-id"
	intersectionName := "Integration Test Intersection"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Step 1: Create intersection
	createRequest := model.CreateIntersectionRequest{
		Name: intersectionName,
		Details: struct {
			Address  string `json:"address"  example:"Corner of Foo and Bar"`
			City     string `json:"city"     example:"Pretoria"`
			Province string `json:"province" example:"Gauteng"`
		}{
			Address:  "123 Integration Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "tjunction",
			Green:            10,
			Yellow:           3,
			Red:              7,
			Speed:            60,
			Seed:             12345,
		},
	}

	createdIntersectionID := "integration-intersection-id"

	expectedIntersection := model.Intersection{
		Name: intersectionName,
		Details: model.Details{
			Address:  "123 Integration Street",
			City:     "Cape Town",
			Province: "Western Cape",
		},
		TrafficDensity: "high",
		DefaultParameters: model.OptimisationParameters{
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}

	suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
		Return(&intersectionpb.IntersectionResponse{Id: createdIntersectionID}, nil)

	suite.userClient.On("AddIntersectionID", ctx, userID, createdIntersectionID).
		Return(nil, nil)

	createResult, err := suite.service.CreateIntersection(ctx, userID, createRequest)
	suite.Require().NoError(err)
	suite.Equal(createdIntersectionID, createResult.Id)

	// Step 2: Get intersection (after creation, user should have access)
	expectedIntersectionIDs := []string{createdIntersectionID}

	getIntersectionResponse := createTestIntersection(
		createdIntersectionID,
		intersectionName,
		"123 Integration Street",
		"Cape Town",
		"Western Cape",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		0,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

	// Mock user intersection IDs stream for GetIntersectionByID
	mockUserStreamGet := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStreamGet.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStreamGet.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).Return(mockUserStreamGet, nil).Once()
	suite.intrClient.On("GetIntersection", ctx, createdIntersectionID).
		Return(getIntersectionResponse, nil)

	getResult, err := suite.service.GetIntersectionByID(ctx, userID, createdIntersectionID)
	suite.Require().NoError(err)
	suite.Equal(createdIntersectionID, getResult.ID)
	suite.Equal(intersectionName, getResult.Name)

	// Step 3: Update intersection
	updateRequest := model.UpdateIntersectionRequest{
		Name: "Updated Integration Intersection",
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

	updatedIntersection := createTestIntersection(
		createdIntersectionID,
		"Updated Integration Intersection",
		"456 Updated Street",
		"Johannesburg",
		"Gauteng",
		commonpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		0,
		commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	)

	// Mock user intersection IDs stream for UpdateIntersectionByID
	mockUserStreamUpdate := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStreamUpdate.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStreamUpdate.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(mockUserStreamUpdate, nil).
		Once()
	suite.intrClient.On("UpdateIntersection", ctx, createdIntersectionID, "Updated Integration Intersection", model.Details{
		Address:  "456 Updated Street",
		City:     "Johannesburg",
		Province: "Gauteng",
	}).
		Return(updatedIntersection, nil)

	updateResult, err := suite.service.UpdateIntersectionByID(
		ctx,
		userID,
		createdIntersectionID,
		updateRequest,
	)
	suite.Require().NoError(err)
	suite.Equal("Updated Integration Intersection", updateResult.Name)
	suite.Equal("456 Updated Street", updateResult.Details.Address)

	// Step 4: Optimise intersection
	mockUserStreamOptimise := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStreamOptimise.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStreamOptimise.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(mockUserStreamOptimise, nil).
		Once()

	// Mock the intersection retrieval for optimisation
	suite.intrClient.On("GetIntersection", ctx, createdIntersectionID).
		Return(getIntersectionResponse, nil)

	// Mock the optimisation service call
	suite.optiClient.On("RunOptimisation", ctx, mock.AnythingOfType("model.OptimisationParameters")).
		Return(&optimisationpb.OptimisationParameters{
			OptimisationType: optimisationpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &commonpb.SimulationParameters{
				IntersectionType: optimisationpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            12,
				Yellow:           3,
				Red:              5,
				Speed:            60,
				Seed:             12345,
			},
		}, nil)

	// Mock the put optimisation call
	suite.intrClient.On("PutOptimisation", ctx, createdIntersectionID, mock.AnythingOfType("model.OptimisationParameters")).
		Return(&intersectionpb.PutOptimisationResponse{Improved: true}, nil)

	err = suite.service.OptimiseIntersectionByID(ctx, userID, createdIntersectionID)
	suite.Require().NoError(err)

	// Step 5: Delete intersection
	mockUserStreamDelete := suite.NewMockUserIntersectionIDsStream()
	for _, id := range expectedIntersectionIDs {
		mockUserStreamDelete.On("Recv").
			Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
			Once()
	}
	mockUserStreamDelete.On("Recv").Return(nil, io.EOF).Once()

	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(mockUserStreamDelete, nil).
		Once()
	suite.userClient.On("RemoveIntersectionID", ctx, userID, createdIntersectionID).Return(nil, nil)
	suite.intrClient.On("DeleteIntersection", ctx, createdIntersectionID).Return(nil, nil)

	err = suite.service.DeleteIntersectionByID(ctx, userID, createdIntersectionID)
	suite.Require().NoError(err)

	// Assert all expectations
	suite.intrClient.AssertExpectations(suite.T())
	suite.userClient.AssertExpectations(suite.T())
	suite.optiClient.AssertExpectations(suite.T())
	mockUserStreamGet.AssertExpectations(suite.T())
	mockUserStreamUpdate.AssertExpectations(suite.T())
	mockUserStreamOptimise.AssertExpectations(suite.T())
	mockUserStreamDelete.AssertExpectations(suite.T())
}

func (suite *TestIntegrationSuite) TestGetAllIntersectionsWithMultipleIntersections() {
	userID := "multi-user-id"

	expectedIntersections := []*intersectionpb.IntersectionResponse{
		createTestIntersection(
			"intersection-1",
			"First Intersection",
			"123 First St",
			"Pretoria",
			"Gauteng",
			commonpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISED,
			5,
			commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
		),
		createTestIntersection(
			"intersection-2",
			"Second Intersection",
			"456 Second Ave",
			"Cape Town",
			"Western Cape",
			commonpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
			0,
			commonpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM,
		),
		createTestIntersection(
			"intersection-3",
			"Third Intersection",
			"789 Third Rd",
			"Johannesburg",
			"Gauteng",
			commonpb.IntersectionStatus_INTERSECTION_STATUS_OPTIMISING,
			2,
			commonpb.TrafficDensity_TRAFFIC_DENSITY_LOW,
		),
	}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	mockStream := suite.NewMockIntersectionStream()

	// Setup stream mock to return intersections one by one
	for _, intersection := range expectedIntersections {
		mockStream.On("Recv").Return(intersection, nil).Once()
	}
	mockStream.On("Recv").Return(nil, io.EOF).Once()

	suite.intrClient.On("GetAllIntersections", ctx).Return(mockStream, nil)

	result, err := suite.service.GetAllIntersections(ctx, userID)

	suite.Require().NoError(err)
	suite.Len(result.Intersections, 3)

	// Verify all intersections
	suite.Equal("intersection-1", result.Intersections[0].ID)
	suite.Equal("First Intersection", result.Intersections[0].Name)
	suite.Equal("INTERSECTION_STATUS_OPTIMISED", result.Intersections[0].Status)

	suite.Equal("intersection-2", result.Intersections[1].ID)
	suite.Equal("Second Intersection", result.Intersections[1].Name)
	suite.Equal("INTERSECTION_STATUS_UNOPTIMISED", result.Intersections[1].Status)

	suite.Equal("intersection-3", result.Intersections[2].ID)
	suite.Equal("Third Intersection", result.Intersections[2].Name)
	suite.Equal("INTERSECTION_STATUS_OPTIMISING", result.Intersections[2].Status)

	suite.intrClient.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}

func (suite *TestIntegrationSuite) TestCascadingErrorScenarios() {
	userID := "error-user-id"
	intersectionID := "error-intersection-id"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Test 1: User service down affects all operations
	suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
		Return(nil, errs.NewUnavailableError("user service is down", map[string]any{})).
		Times(4) // Will be called for get, update, delete, and optimise

	// Test GetIntersectionByID with user service down
	_, err := suite.service.GetIntersectionByID(ctx, userID, intersectionID)
	suite.Require().Error(err)
	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)

	// Test UpdateIntersectionByID with user service down
	updateRequest := model.UpdateIntersectionRequest{Name: "Test"}
	_, err = suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, updateRequest)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)

	// Test DeleteIntersectionByID with user service down
	err = suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)

	// Test OptimiseIntersectionByID with user service down
	err = suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnavailable, svcError.Code)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestIntegrationSuite) TestPermissionConsistencyAcrossOperations() {
	userID := "permission-user-id"
	intersectionID := "permission-intersection-id"

	// User only has access to different intersections
	userIntersectionIDs := []string{"other-intersection-1", "other-intersection-2"}

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Create 4 separate mock streams for each operation
	for i := 0; i < 4; i++ {
		mockUserStream := suite.NewMockUserIntersectionIDsStream()
		for _, id := range userIntersectionIDs {
			mockUserStream.On("Recv").
				Return(&userpb.IntersectionIDResponse{IntersectionId: id}, nil).
				Once()
		}
		mockUserStream.On("Recv").Return(nil, io.EOF).Once()
		suite.userClient.On("GetUserIntersectionIDs", ctx, userID).
			Return(mockUserStream, nil).
			Once()
	}

	// Test GetIntersectionByID - should be forbidden
	_, err := suite.service.GetIntersectionByID(ctx, userID, intersectionID)
	suite.Require().Error(err)
	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test UpdateIntersectionByID - should be forbidden
	updateRequest := model.UpdateIntersectionRequest{Name: "Test"}
	_, err = suite.service.UpdateIntersectionByID(ctx, userID, intersectionID, updateRequest)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test DeleteIntersectionByID - should be forbidden
	err = suite.service.DeleteIntersectionByID(ctx, userID, intersectionID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	// Test OptimiseIntersectionByID - should be forbidden
	err = suite.service.OptimiseIntersectionByID(ctx, userID, intersectionID)
	suite.Require().Error(err)
	svcError, ok = err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrForbidden, svcError.Code)

	suite.userClient.AssertExpectations(suite.T())
}

func (suite *TestIntegrationSuite) TestContextPropagation() {
	userID := "context-user-id"

	// Test that context is properly propagated through all service calls
	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	// Test GetAllIntersections
	mockStream := suite.NewMockIntersectionStream()
	mockStream.On("Recv").Return(nil, io.EOF)

	// The key test here is that the exact context is passed through
	suite.intrClient.On("GetAllIntersections", ctx).Return(mockStream, nil)

	_, err := suite.service.GetAllIntersections(ctx, userID)
	suite.Require().NoError(err)

	suite.intrClient.AssertExpectations(suite.T())
	mockStream.AssertExpectations(suite.T())
}

func (suite *TestIntegrationSuite) TestCreateIntersectionWithDifferentTrafficDensities() {
	userID := "traffic-user-id"

	logger := slog.Default()
	ctx := middleware.SetLogger(context.Background(), logger)

	testCases := []struct {
		name            string
		trafficDensity  string
		expectedPbValue commonpb.TrafficDensity
	}{
		{
			name:            "Low Traffic",
			trafficDensity:  "low",
			expectedPbValue: commonpb.TrafficDensity_TRAFFIC_DENSITY_LOW,
		},
		{
			name:            "Medium Traffic",
			trafficDensity:  "medium",
			expectedPbValue: commonpb.TrafficDensity_TRAFFIC_DENSITY_MEDIUM,
		},
		{
			name:            "High Traffic",
			trafficDensity:  "high",
			expectedPbValue: commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
		},
	}

	for i, tc := range testCases {
		intersectionID := "test-intersection-id"

		createRequest := model.CreateIntersectionRequest{
			Name: tc.name,
			Details: struct {
				Address  string `json:"address"  example:"Corner of Foo and Bar"`
				City     string `json:"city"     example:"Pretoria"`
				Province string `json:"province" example:"Gauteng"`
			}{
				Address:  "123 Test Street",
				City:     "Test City",
				Province: "Test Province",
			},
			TrafficDensity: tc.trafficDensity,
			DefaultParameters: model.SimulationParameters{
				IntersectionType: "tjunction",
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		}

		expectedIntersection := model.Intersection{
			Name: tc.name,
			Details: model.Details{
				Address:  "123 Test Street",
				City:     "Test City",
				Province: "Test Province",
			},
			TrafficDensity: tc.trafficDensity,
			DefaultParameters: model.OptimisationParameters{
				SimulationParameters: model.SimulationParameters{
					IntersectionType: "tjunction",
					Green:            10,
					Yellow:           3,
					Red:              7,
					Speed:            60,
					Seed:             12345,
				},
			},
		}

		suite.intrClient.On("CreateIntersection", ctx, expectedIntersection).
			Return(&intersectionpb.IntersectionResponse{Id: intersectionID}, nil).Once()

		suite.userClient.On("AddIntersectionID", ctx, userID, intersectionID).
			Return(nil, nil).Once()

		_, err := suite.service.CreateIntersection(ctx, userID, createRequest)
		suite.Require().NoError(err, "Failed for test case %d: %s", i, tc.name)
	}

	suite.intrClient.AssertExpectations(suite.T())
	suite.userClient.AssertExpectations(suite.T())
}
