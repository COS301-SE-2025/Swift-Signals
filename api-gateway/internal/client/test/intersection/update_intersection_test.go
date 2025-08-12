package intersection

import (
	"context"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (suite *TestSuite) TestUpdateIntersection_Success() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	name := "Updated Intersection"
	details := model.Details{
		Address:  "456 Updated St",
		City:     "Updated City",
		Province: "Updated Province",
	}

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   intersectionID,
		Name: name,
		Details: &intersectionpb.IntersectionDetails{
			Address:  details.Address,
			City:     details.City,
			Province: details.Province,
		},
		CreatedAt:      timestamppb.Now(),
		LastRunAt:      timestamppb.Now(),
		Status:         intersectionpb.IntersectionStatus_INTERSECTION_STATUS_UNOPTIMISED,
		RunCount:       5,
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	suite.grpcClient.On("UpdateIntersection",
		mock.MatchedBy(func(ctx context.Context) bool {
			deadline, hasDeadline := ctx.Deadline()
			if !hasDeadline {
				return false
			}
			timeUntilDeadline := time.Until(deadline)
			return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
		}),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID &&
				req.Name == name &&
				req.Details.Address == details.Address &&
				req.Details.City == details.City &&
				req.Details.Province == details.Province
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_NotFound() {
	// Arrange
	ctx := context.Background()
	intersectionID := "non-existent-id"
	name := "Updated Intersection"
	details := model.Details{
		Address:  "456 Updated St",
		City:     "Updated City",
		Province: "Updated Province",
	}

	grpcErr := status.Error(codes.NotFound, "intersection not found")

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_EmptyID() {
	// Arrange
	ctx := context.Background()
	intersectionID := ""
	name := "Updated Intersection"
	details := model.Details{
		Address:  "456 Updated St",
		City:     "Updated City",
		Province: "Updated Province",
	}

	grpcErr := status.Error(codes.InvalidArgument, "intersection id is required")

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == ""
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_EmptyName() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	name := ""
	details := model.Details{
		Address:  "456 Updated St",
		City:     "Updated City",
		Province: "Updated Province",
	}

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   intersectionID,
		Name: "",
		Details: &intersectionpb.IntersectionDetails{
			Address:  details.Address,
			City:     details.City,
			Province: details.Province,
		},
	}

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID &&
				req.Name == "" &&
				req.Details.Address == details.Address &&
				req.Details.City == details.City &&
				req.Details.Province == details.Province
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_EmptyDetails() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	name := "Updated Intersection"
	details := model.Details{
		Address:  "",
		City:     "",
		Province: "",
	}

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   intersectionID,
		Name: name,
		Details: &intersectionpb.IntersectionDetails{
			Address:  "",
			City:     "",
			Province: "",
		},
	}

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID &&
				req.Name == name &&
				req.Details.Address == "" &&
				req.Details.City == "" &&
				req.Details.Province == ""
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_Unauthorized() {
	// Arrange
	ctx := context.Background()
	intersectionID := "unauthorized-intersection-id"
	name := "Updated Intersection"
	details := model.Details{
		Address:  "456 Updated St",
		City:     "Updated City",
		Province: "Updated Province",
	}

	grpcErr := status.Error(codes.PermissionDenied, "access denied")

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_InternalError() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	name := "Updated Intersection"
	details := model.Details{
		Address:  "456 Updated St",
		City:     "Updated City",
		Province: "Updated Province",
	}

	grpcErr := status.Error(codes.Internal, "internal server error")

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, grpcErr)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_ContextTimeout() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	intersectionID := "test-intersection-id"
	name := "Updated Intersection"
	details := model.Details{
		Address:  "456 Updated St",
		City:     "Updated City",
		Province: "Updated Province",
	}

	suite.grpcClient.On("UpdateIntersection",
		mock.Anything,
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID
		})).Return(nil, context.DeadlineExceeded)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().Error(err)
	suite.Nil(result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_LongValues() {
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id-with-very-long-name-that-might-cause-issues"
	name := "This is a very long intersection name that might test the limits of the system and ensure proper handling of extended text content"
	details := model.Details{
		Address:  "123456789 Very Long Street Name That Goes On And On And Might Test System Limits",
		City:     "Very Long City Name That Also Tests The System Limits",
		Province: "Very Long Province Name For Testing",
	}

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   intersectionID,
		Name: name,
		Details: &intersectionpb.IntersectionDetails{
			Address:  details.Address,
			City:     details.City,
			Province: details.Province,
		},
	}

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID &&
				req.Name == name &&
				req.Details.Address == details.Address &&
				req.Details.City == details.City &&
				req.Details.Province == details.Province
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_SpecialCharacters() {
	// Test with special characters in the input
	// Arrange
	ctx := context.Background()
	intersectionID := "test-intersection-id"
	name := "Intersection with Special Chars: àáâãäåæçèéêë & @#$%^&*()"
	details := model.Details{
		Address:  "123 Special St. & Ave. #2",
		City:     "São Paulo",
		Province: "Província Especial",
	}

	expectedResponse := &intersectionpb.IntersectionResponse{
		Id:   intersectionID,
		Name: name,
		Details: &intersectionpb.IntersectionDetails{
			Address:  details.Address,
			City:     details.City,
			Province: details.Province,
		},
	}

	suite.grpcClient.On("UpdateIntersection",
		mock.AnythingOfType("*context.timerCtx"),
		mock.MatchedBy(func(req *intersectionpb.UpdateIntersectionRequest) bool {
			return req.Id == intersectionID &&
				req.Name == name &&
				req.Details.Address == details.Address &&
				req.Details.City == details.City &&
				req.Details.Province == details.Province
		})).Return(expectedResponse, nil)

	// Act
	result, err := suite.client.UpdateIntersection(ctx, intersectionID, name, details)

	// Assert
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, result)
	suite.grpcClient.AssertExpectations(suite.T())
}

func TestIntersectionClientUpdateIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
