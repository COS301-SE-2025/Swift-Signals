package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // for development using grpcurl
)

func main() {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collection := client.Database("IntersectionService").Collection("Intersections")
	repo := db.NewMongoIntersectionRepo(collection)
	svc := service.NewIntersectionService(repo)
	h := handler.NewIntersectionHandler(svc)

	lis, err := net.Listen("tcp", ":"+os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer) // for development using grpcurl
	intersectionpb.RegisterIntersectionServiceServer(grpcServer, h)

	log.Println("gRPC server running on :" + os.Getenv("APP_PORT"))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
