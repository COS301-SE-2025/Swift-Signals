package intersection

import (
	"testing"
	"time"

	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/client"
	grpcmocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"
	commonpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/common/v1"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/user/v1"
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
	status commonpb.IntersectionStatus,
	runCount int32,
	trafficDensity commonpb.TrafficDensity,
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
		DefaultParameters: &commonpb.OptimisationParameters{
			OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &commonpb.SimulationParameters{
				IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
		BestParameters: &commonpb.OptimisationParameters{
			OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &commonpb.SimulationParameters{
				IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
				Green:            10,
				Yellow:           3,
				Red:              7,
				Speed:            60,
				Seed:             12345,
			},
		},
		CurrentParameters: &commonpb.OptimisationParameters{
			OptimisationType: commonpb.OptimisationType_OPTIMISATION_TYPE_GRIDSEARCH,
			Parameters: &commonpb.SimulationParameters{
				IntersectionType: commonpb.IntersectionType_INTERSECTION_TYPE_TJUNCTION,
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
