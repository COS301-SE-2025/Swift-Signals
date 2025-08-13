package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/client"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/handler"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/service"
	_ "github.com/COS301-SE-2025/Swift-Signals/api-gateway/swagger"
	"github.com/COS301-SE-2025/Swift-Signals/shared/config"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Port             int    `env:"PORT"           envDefault:"9090"`
	JwtSecret        string `env:"JWT_SECRET"     envDefault:"jwt-test-secret"`
	UserServiceAddr  string `env:"USER_GRPC_ADDR" envDefault:"localhost:50051"` // TODO: Change to proper address
	IntersectionAddr string `env:"INTR_GRPC_ADDR" envDefault:"localhost:50052"` // TODO: Change to proper address
	OptimisationAddr string `env:"OPTI_GRPC_ADDR" envDefault:"localhost:50054"` // TODO: Change to proper address
}

// @title Authentication API Gateway
// @version 1.0
// @description This is the API Gateway for the Swift-Signals project, forwarding requests to backend gRPC microservices.
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
	var cfg Config
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	userClient := mustConnectUserService(cfg.UserServiceAddr)
	intrClient := mustConnectIntersectionService(cfg.IntersectionAddr)
	optiClient := mustConnectOptimisationService(cfg.OptimisationAddr)

	baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	mux := setupRoutes(baseLogger, cfg.JwtSecret, userClient, intrClient, optiClient)

	server := createServer(cfg.Port, mux)
	runServer(server)
}

func mustConnectUserService(address string) *client.UserClient {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	) // TODO: Add TLS
	if err != nil {
		log.Fatalf("failed to connect to User gRPC server: %v", err)
	}
	log.Println("Connected to User-Service")
	return client.NewUserClientFromConn(conn)
}

func mustConnectIntersectionService(address string) *client.IntersectionClient {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	) // TODO: Add TLS
	if err != nil {
		log.Fatalf("failed to connect to Intersection gRPC server: %v", err)
	}
	log.Println("Connected to Intersection-Service")
	return client.NewIntersectionClient(conn)
}

func mustConnectOptimisationService(address string) *client.OptimisationClient {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	) // TODO: Add TLS
	if err != nil {
		log.Fatalf("failed to connect to Optimisation gRPC server: %v", err)
	}
	log.Println("Connected to Optimisation-Service")
	return client.NewOptimisationClientFromConn(conn)
}

func setupRoutes(
	logger *slog.Logger,
	JwtSecret string,
	userClient *client.UserClient,
	intrClient *client.IntersectionClient,
	optiClient *client.OptimisationClient,
) http.Handler {
	mux := http.NewServeMux()

	// Auth routes
	authService := service.NewAuthService(userClient)
	authHandler := handler.NewAuthHandler(authService)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /logout", authHandler.Logout)
	mux.HandleFunc("POST /reset-password", authHandler.ResetPassword)
	log.Println("Initialized Auth Handlers.")

	// Intersection routes
	intersectionService := service.NewIntersectionService(intrClient, optiClient, userClient)
	intersectionHandler := handler.NewIntersectionHandler(intersectionService)
	mux.HandleFunc("GET /intersections", intersectionHandler.GetAllIntersections)
	mux.HandleFunc("GET /intersections/{id}", intersectionHandler.GetIntersection)
	mux.HandleFunc("POST /intersections", intersectionHandler.CreateIntersection)
	mux.HandleFunc("PATCH /intersections/{id}", intersectionHandler.UpdateIntersection)
	mux.HandleFunc("DELETE /intersections/{id}", intersectionHandler.DeleteIntersection)
	mux.HandleFunc("POST /intersections/{id}/optimise", NotImplemented)
	mux.HandleFunc("GET /intersections/simple", NotImplemented)
	log.Println("Initialized Intersection Handlers.")

	// Swagger
	mux.Handle("/docs/", httpSwagger.WrapHandler)
	log.Println("Swagger UI available at http://localhost:9090/docs/index.html")

	// Middleware stack
	return middleware.CreateStack(
		middleware.Logging(logger),
		middleware.CORS,
		middleware.AuthMiddleware(
			JwtSecret,
			"/login",
			"/register",
			"/reset-password",
			"/docs",
			"/favicon.ico",
		),
	)(
		mux,
	)
}

func createServer(port int, handler http.Handler) *http.Server {
	addr := fmt.Sprintf(":%d", port)
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func runServer(server *http.Server) {
	// Start server in goroutine
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received

	// Gracefully shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "NotImplemented", http.StatusNotImplemented)
}
