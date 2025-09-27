package test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/swiftsignals/intersection/v1"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type IntegrationTestSuite struct {
	suite.Suite
	mongoContainer *mongodb.MongoDBContainer
	mongoClient    *mongo.Client
	grpcServer     *grpc.Server
	bufListener    *bufconn.Listener
	grpcConn       *grpc.ClientConn
	client         intersectionpb.IntersectionServiceClient
	handler        *handler.Handler
	service        service.IntersectionService
	repo           db.IntersectionRepository
	ctx            context.Context
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	mongoContainer, err := mongodb.Run(suite.ctx,
		"mongo:6.0",
	)
	suite.Require().NoError(err)
	suite.mongoContainer = mongoContainer

	mongoURI, err := mongoContainer.ConnectionString(suite.ctx)
	suite.Require().NoError(err)

	client, err := mongo.Connect(suite.ctx, options.Client().ApplyURI(mongoURI))
	suite.Require().NoError(err)
	suite.mongoClient = client

	collection := client.Database("IntersectionService").Collection("Intersections")
	suite.repo = db.NewMongoIntersectionRepo(collection)
	suite.service = service.NewIntersectionService(suite.repo)
	suite.handler = handler.NewIntersectionHandler(suite.service)

	suite.bufListener = bufconn.Listen(1024 * 1024)
	suite.grpcServer = grpc.NewServer()
	intersectionpb.RegisterIntersectionServiceServer(suite.grpcServer, suite.handler)

	go func() {
		if err := suite.grpcServer.Serve(suite.bufListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return suite.bufListener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	suite.Require().NoError(err)
	suite.grpcConn = conn
	suite.client = intersectionpb.NewIntersectionServiceClient(conn)
}

func (suite *IntegrationTestSuite) SetupTest() {
	err := suite.mongoClient.Database("IntersectionService").
		Collection("Intersections").
		Drop(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.grpcConn != nil {
		if err := suite.grpcConn.Close(); err != nil {
			suite.T().Logf("Failed to close gRPC connection: %v", err)
		}
	}
	if suite.grpcServer != nil {
		suite.grpcServer.Stop()
	}
	if suite.mongoClient != nil {
		if err := suite.mongoClient.Disconnect(suite.ctx); err != nil {
			suite.T().Logf("Failed to disconnect MongoDB client: %v", err)
		}
	}
	if suite.mongoContainer != nil {
		if err := suite.mongoContainer.Terminate(suite.ctx); err != nil {
			suite.T().Logf("Failed to terminate MongoDB container: %v", err)
		}
	}
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
