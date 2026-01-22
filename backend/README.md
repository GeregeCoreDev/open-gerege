# Template Backend API

![CI](https://github.com/gerege-core/backend-refactor-v25/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/git.gerege.mn/backend-packages/template-v25)](https://goreportcard.com/report/git.gerege.mn/backend-packages/template-v25)

Go backend API template built with Clean Architecture principles using Fiber v2, GORM, and PostgreSQL. Now world-class with automated CI/CD and observability.

## Features

- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Fiber v2**: High-performance web framework
- **GORM**: ORM with PostgreSQL support
- **SSO Integration**: Single Sign-On with session caching
- **RBAC**: Role-Based Access Control
- **Observability**: OpenTelemetry tracing & Prometheus metrics
- **Integration Testing**: Ephemeral databases using Testcontainers
- **CI/CD**: GitHub Actions pipeline with strict linting & security checks
- **Swagger**: Auto-generated API documentation
- **Structured Logging**: Zap logger with request ID propagation
- **Security Headers**: CSP, CORS, rate limiting
- **Graceful Shutdown**: Clean resource cleanup

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/                  # Dependency injection container
│   ├── auth/                 # Authentication middleware
│   ├── db/                   # Database connection
│   ├── domain/               # Domain models/entities
│   ├── http/
│   │   ├── dto/              # Data transfer objects
│   │   ├── handlers/         # HTTP handlers
│   │   └── router/           # Route definitions
│   ├── middleware/           # HTTP middlewares
│   ├── repository/           # Data access layer
│   └── service/              # Business logic layer
├── docs/                     # Swagger generated docs
├── .github/workflows/        # CI/CD pipelines
└── docker/                   # Docker configurations
```

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 15+
- Make (optional)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd golang-template-v25-main

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run the server
go run cmd/server/main.go
```

### Using Make

```bash
make run        # Run the server
make build      # Build binary
make test       # Run unit tests
make test-integration # Run integration tests (requires Docker)
make test-all   # Run all tests
make audit      # Run security audit (govulncheck) and linter
make mocks      # Generate mocks (requires mockery)
make lint       # Run linter
make swagger    # Generate Swagger docs
```

## Configuration

Environment variables (`.env` file):

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8000
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=template_db

# Auth
AUTH_CACHE_TTL=1h
AUTH_CACHE_MAX=10000

# TLS (optional)
TLS_CERT=
TLS_KEY=
```

## API Endpoints

### Public Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check with DB status |
| GET | `/docs/*` | Swagger UI |

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| GET | `/auth/login` | SSO redirect |
| GET | `/auth/callback` | OAuth2 callback |
| POST | `/auth/logout` | Logout |
| GET | `/auth/verify` | Token verification |

### Protected Routes (require authentication)

| Resource | Endpoints |
|----------|-----------|
| User | `/user/*` |
| Role | `/role/*` |
| Permission | `/permission/*` |
| Organization | `/organization/*` |
| System | `/system/*` |
| Module | `/module/*` |

## Health Check

The `/health` endpoint returns comprehensive status information:

```json
{
  "code": "OK",
  "data": {
    "status": "ok",
    "uptime": 3600,
    "timestamp": "2025-01-10T12:00:00Z",
    "database": {
      "status": "ok",
      "open_conns": 10,
      "in_use": 2,
      "idle": 8
    }
  }
}
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Generating Swagger Docs

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/server/main.go -o docs
```

## Deployment

### Docker

```bash
# Build image
docker build -t template-backend .

# Run container
docker run -p 8000:8000 --env-file .env template-backend
```

### Docker Compose

```bash
docker-compose up -d
```

## CI/CD

GitHub Actions workflows are configured for:

- **Lint**: golangci-lint checks
- **Test**: Unit tests with PostgreSQL service
- **Build**: Binary compilation
- **Docker**: Image build and push
- **Security**: Gosec and Trivy scans

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is proprietary software of Gerege Core Team.

## Authors

- Bayarsaikhan Otgonbayar, CTO - Gerege Core Team
