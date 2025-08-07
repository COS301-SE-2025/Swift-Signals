package main

import (
	// Un/comment for Postgresql
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // for development using grpcurl
)

func main() {
	// Postgresql Connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Failed to close DB connection: %v", err)
		}
	}()

	repo := db.NewPostgresUserRepo(dbConn)

	svc := service.NewUserService(repo)
	handler := handler.NewUserHandler(svc)

	lis, err := net.Listen("tcp", ":"+os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer) // for development using grpcurl

	userpb.RegisterUserServiceServer(grpcServer, handler)

	log.Println("gRPC server running on :" + os.Getenv("APP_PORT"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("gRPC server exited with error: %v", err)
	}
}
