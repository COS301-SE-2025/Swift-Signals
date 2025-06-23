package main

import (
	"log"
	"net"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" //for development using grpcurl
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer) //for development using grpcurl

	userRepo := db.NewUserRepository()
	userService := service.NewService(userRepo)
	userHandler := handler.NewHandler(userService)

	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Println("gRPC server listening on :50051")
	grpcServer.Serve(lis)
}
