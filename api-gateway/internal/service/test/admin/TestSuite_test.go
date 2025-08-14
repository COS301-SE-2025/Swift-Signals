package admin

import (
	mocks "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	client  *mocks.MockUserClientInterface
	service service.AdminServiceInterface
}

func (suite *TestSuite) SetupTest() {
	suite.client = new(mocks.MockUserClientInterface)
	suite.service = service.NewAdminService(suite.client)
}

// createTestUser creates a test user protobuf object
func createTestUser(
	id, name, email string,
	isAdmin bool,
	intersectionIDs []string,
) *userpb.UserResponse {
	return &userpb.UserResponse{
		Id:              id,
		Name:            name,
		Email:           email,
		IsAdmin:         isAdmin,
		IntersectionIds: intersectionIDs,
	}
}
