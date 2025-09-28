package test

import (
	"context"
	"time"

	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
)

func (suite *IntegrationTestSuite) TestGetAllIntersections() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	createReq2 := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test2 Intersection",
		TrafficDensity: commonpb.TrafficDensity_TRAFFIC_DENSITY_LOW,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	_, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	_, err = suite.client.CreateIntersection(ctx, createReq2)
	suite.Require().NoError(err)

	req := &intersectionpb.GetAllIntersectionsRequest{
		Page:     1,
		PageSize: 2,
	}

	resp, err := suite.client.GetAllIntersections(ctx, req)

	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}
