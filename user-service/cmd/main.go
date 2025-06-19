package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal"
	userpb "github.com/COS301-SE-2025/Swift-Signals/user-service/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" //for development using grpcurl
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

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
	defer dbConn.Close()

	repo := db.NewPostgresUserRepo(dbConn)
	svc := user.NewService(repo)
	handler := user.NewHandler(svc)

	lis, err := net.Listen("tcp", ":"+os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer) //for development using grpcurl

	userpb.RegisterUserServiceServer(grpcServer, handler)

	log.Println("gRPC server running on :" + os.Getenv("APP_PORT"))
	grpcServer.Serve(lis)
}
