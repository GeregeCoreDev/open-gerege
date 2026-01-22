# Open-Gerege

Enterprise-grade backend API template and web application framework built with Go and Next.js.

## Tech Stack

### Backend
- **Go 1.25** with Fiber v2 framework
- **PostgreSQL** with GORM ORM
- **Redis** for session/cache management
- **OpenTelemetry** for observability
- **Swagger** for API documentation

### Frontend
- **Next.js 16** with React 19
- **TypeScript 5**
- **Tailwind CSS 4**
- **Zustand** for state management

## Features

- SSO integration + local authentication with MFA/TOTP
- Role-Based Access Control (RBAC)
- User and organization management
- Notification system
- API request logging and analytics
- Health monitoring with Prometheus metrics
- Clean Architecture with DDD patterns

## Project Structure

```
open-gerege/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── domain/          # Domain entities
│   │   ├── service/         # Business logic
│   │   ├── repository/      # Data access layer
│   │   └── http/            # HTTP handlers & routes
│   ├── migrations/          # Database migrations
│   └── docs/                # Swagger documentation
└── frontend/
    └── src/
        ├── app/             # Next.js App Router
        ├── components/      # UI components
        ├── features/        # Feature modules
        └── lib/             # Utilities
```

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 20+
- PostgreSQL 15+
- Redis (optional, for session caching)

### Backend Setup

```bash
cd backend

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Run database migrations
make migrate

# Start server
make run
```

The backend server will start at `http://localhost:8000`

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start at `http://localhost:2000`

## Development

### Backend Commands

```bash
make run              # Run with live reload
make build            # Build binary
make test             # Run tests
make test-integration # Run integration tests
make swagger          # Generate Swagger docs
make lint             # Run linter
make audit            # Security audit
```

### Frontend Commands

```bash
npm run dev           # Development server with Turbopack
npm run build         # Production build
npm run start         # Start production server
npm run lint          # ESLint
```

## API Documentation

Swagger documentation is available at `http://localhost:8000/swagger/index.html` when the server is running.

## Docker

```bash
# Development
docker-compose up

# Production build
docker build -t open-gerege-backend ./backend
```

## License

MIT License
