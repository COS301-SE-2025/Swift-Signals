package test

import (
	"context"
	"time"

	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
)

func (suite *IntegrationTestSuite) TestCreateIntersection() {
	req := &intersectionpb.CreateIntersectionRequest{
		Name: "Test Intersection",
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	resp, err := suite.client.CreateIntersection(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

func (suite *IntegrationTestSuite) TestCreateIntersection_Failure() {
	req := &intersectionpb.CreateIntersectionRequest{
		Name: "",
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 30*time.Second)
	defer cancel()

	resp, err := suite.client.CreateIntersection(ctx, req)
	suite.Require().Error(err)
	suite.Require().Nil(resp)
}
