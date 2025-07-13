package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"

	_ "github.com/COS301-SE-2025/Swift-Signals/api-gateway/swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Authentication API Gateway
// @version 1.0
// @description This is the API Gateway for the Swift-Signals project,
// @description forwarding requests to backend gRPC microservices.
// @termsOfService http://example.com/terms/

// @contact.name Inside Insights Team
// @contact.url https://swagger.io/
// @contact.email insideinsights2025@gmail.com

// @host localhost:9090
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // NOTE: Will change to use TLS later on
	if err != nil {
		log.Fatalf("failed to connect to User gRPC server: %v", err)
	}
	userClient := client.NewUserClient(userConn)
	log.Println("Connected to User-Service")

	intrConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure()) // NOTE: Will change to use TLS later on
	if err != nil {
		log.Fatalf("failed to connect to Intersection gRPC server: %v", err)
	}
	intrClient := client.NewIntersectionClient(intrConn)
	log.Println("Connected to Intersection-Service")

	mux := http.NewServeMux()

	authService := service.NewAuthService(userClient)
	authHandler := handler.NewAuthHandler(authService)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /logout", authHandler.Logout)
	mux.HandleFunc("POST /reset-password", authHandler.ResetPassword)
	log.Println("Initialized Auth Handlers.")

	intersectionService := service.NewIntersectionService(intrClient)
	intersectionHandler := handler.NewIntersectionHandler(intersectionService)
	mux.HandleFunc("GET /intersections", intersectionHandler.GetAllIntersections)
	// mux.HandleFunc("GET /intersections/simple", nil)
	mux.HandleFunc("GET /intersections/{id}", intersectionHandler.GetIntersection)
	mux.HandleFunc("POST /intersections", intersectionHandler.CreateIntersection)
	mux.HandleFunc("PATCH /intersections/{id}", intersectionHandler.UpdateIntersection)
	// mux.HandleFunc("DELETE /intersections/{id}", nil)
	// mux.HandleFunc("POST /intersections/{id}/optimise", nil)
	log.Println("Initialized Intersection Handlers.")

	log.Println("Registered API routes.")

	mux.Handle("/docs/", httpSwagger.WrapHandler)
	log.Println("Swagger UI available at http://localhost:9090/docs/index.html")

	serverAddr := fmt.Sprintf(":%d", 9090)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      middleware.CORS(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Block until a signal is received
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
