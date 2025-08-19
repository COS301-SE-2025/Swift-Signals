package test

import (
	"context"
	"time"

	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
)

func (suite *IntegrationTestSuite) TestGetIntersection() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	intersection, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.IntersectionIDRequest{
		Id: intersection.GetId(),
	}

	resp, err := suite.client.GetIntersection(ctx, req)

	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

func (suite *IntegrationTestSuite) TestGetIntersection_Failure() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	_, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.IntersectionIDRequest{
		Id: "invalid id",
	}

	resp, err := suite.client.GetIntersection(ctx, req)

	suite.Require().Error(err)
	suite.Require().Nil(resp)
}
