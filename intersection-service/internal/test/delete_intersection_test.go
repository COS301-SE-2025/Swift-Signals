package test

import (
	"context"
	"time"

	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
)

func (suite *IntegrationTestSuite) TestDeleteIntersection() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name: "Test Intersection",
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	intersection, err := suite.client.CreateIntersection(ctx, createReq)
	if err != nil {
		suite.Require().Error(err)
	}

	req := &intersectionpb.IntersectionIDRequest{
		Id: intersection.GetId(),
	}

	resp, err := suite.client.DeleteIntersection(ctx, req)

	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

func (suite *IntegrationTestSuite) TestDeleteIntersection_Failure() {
	createReq := &intersectionpb.CreateIntersectionRequest{
		Name: "Test Intersection",
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	_, err := suite.client.CreateIntersection(ctx, createReq)
	suite.Require().NoError(err)

	req := &intersectionpb.IntersectionIDRequest{}

	resp, err := suite.client.DeleteIntersection(ctx, req)

	suite.Require().Error(err)
	suite.Require().Nil(resp)
}
