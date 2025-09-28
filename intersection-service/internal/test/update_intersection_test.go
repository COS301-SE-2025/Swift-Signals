package test

import (
	"context"
	"time"

	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
)

func (suite *IntegrationTestSuite) TestUpdateIntersection() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	intersection, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.UpdateIntersectionRequest{
		Id:   intersection.GetId(),
		Name: "New Name",
	}

	resp, err := suite.client.UpdateIntersection(ctx, req)

	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

func (suite *IntegrationTestSuite) TestUpdateIntersection_Failure_Min() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	intersection, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.UpdateIntersectionRequest{
		Id:   intersection.GetId(),
		Name: "N",
	}

	resp, err := suite.client.UpdateIntersection(ctx, req)

	suite.Require().Error(err)
	suite.Require().Nil(resp)
}

func (suite *IntegrationTestSuite) TestUpdateIntersection_Failure_Max() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	intersection, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.UpdateIntersectionRequest{
		Id:   intersection.GetId(),
		Name: "NdfjdkjfkdfjdkfjdkjfdkjfkdjfdkjfkdFDFkjdfjkdfjdkfjdkjfdjfkjdfkdjfkdjfdkjkdjfkdjfkdjfdjfkdfkdjkdkjfdkfjakfjjfaksjfksjdfkasjfksjdfk",
	}

	resp, err := suite.client.UpdateIntersection(ctx, req)

	suite.Require().Error(err)
	suite.Require().Nil(resp)
}

func (suite *IntegrationTestSuite) TestUpdateIntersection_Failure_UUID() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: commonpb.TrafficDensity_TRAFFIC_DENSITY_HIGH,
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	_, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.UpdateIntersectionRequest{
		Id:   "invalid id",
		Name: "New Name",
	}

	resp, err := suite.client.UpdateIntersection(ctx, req)

	suite.Require().Error(err)
	suite.Require().Nil(resp)
}
