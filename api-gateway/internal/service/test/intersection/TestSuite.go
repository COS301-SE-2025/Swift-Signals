package intersection

import (
	"testing"
	"time"

	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/client"
	grpcmocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TestSuite struct {
	suite.Suite
	intrClient *mocks.MockIntersectionClientInterface
	userClient *mocks.MockUserClientInterface
	optiClient *mocks.MockOptimisationClientInterface
	service    service.IntersectionServiceInterface
}

func (suite *TestSuite) SetupTest() {
	suite.intrClient = new(mocks.MockIntersectionClientInterface)
	suite.userClient = new(mocks.MockUserClientInterface)
	suite.optiClient = new(mocks.MockOptimisationClientInterface)
	suite.service = service.NewIntersectionService(
		suite.intrClient,
		suite.optiClient,
		suite.userClient,
	)
}

// Helper method to create mock streams for testing
func (suite *TestSuite) NewMockIntersectionStream() *grpcmocks.MockIntersectionService_GetAllIntersectionsClient[intersectionpb.IntersectionResponse] {
	return grpcmocks.NewMockIntersectionService_GetAllIntersectionsClient[intersectionpb.IntersectionResponse](
		suite.T(),
	)
}

func (suite *TestSuite) NewMockUserIntersectionIDsStream() *grpcmocks.MockUserService_GetUserIntersectionIDsClient[userpb.IntersectionIDResponse] {
	return grpcmocks.NewMockUserService_GetUserIntersectionIDsClient[userpb.IntersectionIDResponse](
		suite.T(),
	)
}

// createTestIntersection creates a test intersection protobuf object
func createTestIntersection(
	id, name, address, city, province string,
	status intersectionpb.IntersectionStatus,
	runCount int32,
	trafficDensity intersectionpb.TrafficDensity,
) *intersectionpb.IntersectionResponse {
	now := time.Now()
	return &intersectionpb.IntersectionResponse{
		Id:   id,
		Name: name,
		Details: &intersectionpb.IntersectionDetails{
			Address:  address,
			City:     city,
			Province: province,
		},
		CreatedAt:      timestamppb.New(now),
		LastRunAt:      timestamppb.New(now),
		Status:         status,
		RunCount:       runCount,
		TrafficDensity: trafficDensity,
		DefaultParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
		BestParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
		CurrentParameters: &intersectionpb.OptimisationParameters{
			OptimisationType: intersectionpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &intersectionpb.SimulationParameters{
				IntersectionType: intersectionpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
	}
}

func TestService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
