package test

import (
	"context"
	"time"

	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
)

func (suite *IntegrationTestSuite) TestPutOptimisation() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	intersection, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.PutOptimisationRequest{
		Id: intersection.GetId(),
	}

	resp, err := suite.client.PutOptimisation(ctx, req)

	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

func (suite *IntegrationTestSuite) TestPutOptimisation_Failure() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: intersectionpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	_, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.PutOptimisationRequest{
		Id: "invalid id",
	}

	resp, err := suite.client.PutOptimisation(ctx, req)

	suite.Require().Error(err)
	suite.Require().Nil(resp)
}
