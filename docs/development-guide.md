## Layout 

```/Swift-Signals
│
├── /configs
│   ├── /development
│   │   ├── services.yaml        # All services config for dev
│   │   └── database.yaml        # Database configs for dev
│   ├── /staging
│   │   ├── services.yaml
│   │   └── database.yaml
│   └── /production
│       ├── services.yaml
│       └── database.yaml
│
├── /scripts
│   ├── setup.sh                 # Dev setup script
│   ├── test.sh                  # Integration test runner
│   ├── seed-data.sh             # Populate test data
│   └── deploy.sh                # Deployment automation
│
├── /build
│   ├── /package
│   │   ├── Dockerfile.user
│   │   ├── Dockerfile.sim 
│   │   ├── Dockerfile.api      # api-gateway
│   │   ├── Dockerfile.control 
│   │   ├── Dockerfile.metrics 
│   │   ├── Dockerfile.ai
│   │   └── Dockerfile.frontend
│   └── /ci
│       └── github-actions.yml   # Primary CI/CD pipeline
│
├── /deployments
│   ├── docker-compose.yml       # Local development setup
│   ├── docker-compose.prod.yml  # Production compose setup
│   └── /k8s
│       ├── /development
│       │   ├── user-service.yaml
│       │   ├── simulation-service.yaml
│       │   ├── api-gateway.yaml
│       │   ├── control-service.yaml
│       │   ├── metrics-service.yaml
│       │   ├── ai-service.yaml
│       │   ├── frontend-service.yaml
│       │   └── namespace.yaml
│       └── /production
│           ├── user-service.yaml
│           ├── simulation-service.yaml
│           ├── api-gateway.yaml
│           ├── control-service.yaml
│           ├── metrics-service.yaml
│           ├── ai-service.yaml
│           ├── frontend-service.yaml
│           └── namespace.yaml
│
├── /test
│   ├── /unit                    # Unit tests
│   ├── /integration             # Integration tests
│   ├── /performance             # Load/stress tests
│   ├── /fixtures                # Test data
│   └── /mocks                   # Mock implementations
│
├── /docs
│   ├── architecture.md
│   ├── api-usage.md
│   ├── development-guide.md
│   ├── /diagrams
│   │   ├── sequence-diagram.png
│   │   └── architecture-overview.svg
│   └── /api
│       ├── swagger.yaml         # OpenAPI/Swagger specs
│       └── /examples            # Request/response examples
│
├── /tools
│   ├── /gen-openapi
│   │   └── generate.sh
│   └── /db-migration
│       └── migrate.sh
│
├── /assets
│   └── /images
│       ├── sequence-diagram.png
│       └── ml-accuracy-plot.jpg
│
├── /api-gateway
│   ├── /cmd
│   │   └── main.go
│   ├── /internal
│   │   ├── /api
│   │   ├── /service
│   │   └── /config
│   ├── README.md               # Service-specific documentation
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile             # Service-specific Dockerfile
│
├── /user-service
│   ├── /cmd
│   │   └── main.go
│   ├── /internal
│   │   ├── /api
│   │   ├── /service
│   │   └── /config
│   ├── /models
│   ├── /pkg
│   ├── /db
│   │   ├── init.js            # MongoDB initialization script
│   │   └── migrations/        # Database schema migrations
│   ├── README.md              # Service-specific documentation
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile             # Service-specific Dockerfile
│
├── /simulation-service
│   ├── /cmd
│   │   └── main.go
│   ├── /internal
│   │   ├── /api
│   │   ├── /service
│   │   └── /config
│   ├── /sumo-config           # Network & route files
│   ├── README.md              # Service-specific documentation
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
│
├── /ai-service
│   ├── /models                # ML models (saved as .pt, .h5, etc.)
│   ├── /inference             # Python AI inference code
│   ├── /training              # Model training scripts
│   ├── /db
│   │   └── mongo-setup.js     # Optional for model metadata
│   ├── /tests                 # Python tests
│   ├── README.md              # Service-specific documentation
│   ├── requirements.txt       # Python dependencies
│   └── Dockerfile
│
├── /control-service
│   ├── /cmd
│   │   └── main.go
│   ├── /internal
│   │   ├── /api
│   │   ├── /service
│   │   └── /config
│   ├── README.md              # Service-specific documentation
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
│
├── /metrics-service
│   ├── /cmd
│   │   └── main.go
│   ├── /internal
│   │   ├── /api
│   │   ├── /service
│   │   └── /config
│   ├── /db
│   │   └── setup.sql          # PostgreSQL or InfluxDB schema
│   ├── README.md              # Service-specific documentation
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
│
├── /frontend
│   ├── /public
│   ├── /src
│   │   ├── /components
│   │   ├── /services          # Frontend service clients
│   │   └── /pages
│   ├── /tests                 # Frontend tests
│   ├── README.md              # Dashboard documentation
│   ├── package.json
│   └── Dockerfile
│
├── /shared
│   ├── go.mod                 # Separate Go module for shared code
│   ├── go.sum
│   ├── /auth
│   │   └── jwt.go             # JWT encode/decode, token validation
│   ├── /db
│   │   └── mongo.go           # MongoDB connection pool setup
│   ├── /config
│   │   └── loader.go          # YAML/ENV config parsing
│   ├── /logger
│   │   └── logger.go          # Zap or logrus wrapper
│   ├── /proto                 # Protocol buffer definitions
│   │   ├── user.proto
│   │   └── simulation.proto
│   └── /client                # Client libraries for services
│       ├── user_client.go
│       └── simulation_client.go
│
├── README.md                  # Project overview and setup guide
├── .env.example               # Template for environment variables
├── .env.development           # Dev environment settings
├── .env.test                  # Test environment settings
├── .gitignore                 # Version control exclusions
├── .dockerignore              # Docker build exclusions
├── .golangci.yml              # Go linting configuration
└── .pre-commit-config.yaml    # Git hooks for quality checks```
