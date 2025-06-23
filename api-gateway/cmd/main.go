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

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/handler"

	_ "github.com/COS301-SE-2025/Swift-Signals/api-gateway/swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Authentication API Gateway
// @version 1.0
// @description This is the API Gateway for the Swift-Signals project,
// @description forwarding requests to backend gRPC microservices.
// @termsOfService http://example.com/terms/

// @contact.name Inside Insights Team
// @contact.url ...
// @contact.email support@example.com

// @host localhost:9090
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	authHandler := handler.NewAuthHandler()
	log.Println("Initialized Auth Handler.")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /logout", authHandler.Logout)
	mux.HandleFunc("POST /reset-password", authHandler.ResetPassword)
	log.Println("Registered API routes.")

	mux.Handle("/docs/", httpSwagger.WrapHandler)
	log.Println("Swagger UI available at http://localhost:9090/docs/index.html")

	serverAddr := fmt.Sprintf(":%d", 9090)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
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
