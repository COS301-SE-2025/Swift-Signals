package main

import (
	"log"
	"net"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" //for development using grpcurl
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer) //for development using grpcurl

	intersectionRepo := db.NewIntersectionRepository()
	intersectionService := service.NewService(intersectionRepo)
	intersectionHandler := handler.NewHandler(intersectionService)

	intersectionpb.RegisterIntersectionServiceServer(grpcServer, intersectionHandler)

	log.Println("gRPC server listening on :50052")
	grpcServer.Serve(lis)
}
