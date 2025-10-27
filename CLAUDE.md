# Identity Server - Claude Development Reference

## Project Overview
A Go-based identity management server providing user registration, authentication, and JWT token generation. Built with a clean architecture separating business logic, storage, and HTTP handling layers.

## Architecture

### Core Components
- **Main Entry Point**: `app/auth-api/main.go` - Application startup, database connection, and server initialization
- **HTTP Handlers**: `app/auth-api/handlers/` - REST API endpoints and routing
- **Business Logic**: `business/` - Core business operations and validation
- **Database Layer**: `business/store/` - Data access, connection management, and migrations
- **Foundation**: `foundation/` - Shared utilities, middleware, and crypto operations

### Key Features
- User registration and authentication
- JWT token generation and validation
- Password hashing and validation
- Database migrations
- CORS support
- Health check endpoints (liveness/readiness)
- MySQL database integration

## Technology Stack
- **Language**: Go 1.24
- **Database**: MySQL
- **Authentication**: JWT tokens with RSA key pairs
- **HTTP Router**: julienschmidt/httprouter
- **Migrations**: golang-migrate/migrate
- **Deployment**: Docker + Kubernetes/Helm

## Dependencies (go.mod)
```go
require (
    github.com/DATA-DOG/go-sqlmock v1.5.0
    github.com/go-sql-driver/mysql v1.6.0
    github.com/golang-jwt/jwt/v4 v4.2.0
    github.com/golang-migrate/migrate v3.5.4+incompatible
    github.com/google/jsonapi v1.0.0
    github.com/google/uuid v1.3.0
    github.com/joho/godotenv v1.4.0
    github.com/julienschmidt/httprouter v1.3.0
    github.com/rs/cors v1.8.2
    github.com/sirupsen/logrus v1.8.1
    github.com/stretchr/testify v1.7.0
    golang.org/x/crypto v0.0.0-20220131195533-30dcbda58838
)
```

## API Endpoints

### Public Endpoints
- `POST /register` - User registration
- `POST /login` - User authentication and token generation
- `GET /liveness` - Kubernetes liveness probe
- `GET /readiness` - Kubernetes readiness probe

### Protected Endpoints (require JWT)
- `GET /user/home` - User profile access
- `DELETE /admin/delete/:id` - User deletion

## Database Schema

### Users Table (`identity_users`)
```sql
CREATE TABLE identity_users (
    id VARCHAR(64) PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(120),
    password VARCHAR(120),
    company VARCHAR(64),
    post_code VARCHAR(64),
    terms INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
```

## Configuration

### Environment Variables
```bash
# Development (.env)
PORT="8089"
ENV="dev"
MYSQL_DATABASE="identity"
MYSQL_USERNAME="username"
MYSQL_PASSWORD="password"
MYSQL_PORT="3306"
MYSQL_HOST="127.0.0.1"
BASE_PATH="/"


# Additional Production Variables
ISSUER="riya-dennis"
KEY_PATH="/tmp/keys"
```

## Development Workflow

### Running Locally
1. Set up `.env` file with database credentials
2. Ensure MySQL is running and accessible
3. Run migrations automatically on startup
4. Start server: `go run app/auth-api/main.go`

### Testing
- Test files: `**/*_test.go`
- Current coverage: ~28-33%
- Test database required (`.env_test` file)
- Key test areas:
  - Handler tests: `app/auth-api/handlers/*_test.go`
  - Store tests: `business/store/*_test.go`
  - Validation tests: `business/validation/validate_test.go`

### Docker Build & Deploy
```bash
# Build
make docker-build  # docker build -t riyadennis/identity-server:1.3.0 .

# Push
make docker-push   # docker push riyadennis/identity-server:1.3.0

# Run locally
make docker-run    # Run with production config
```

### Kubernetes Deployment
```bash
# Helm deployment
make helm-install    # helm install identity ./zarf/identity
make helm-uninstall  # helm uninstall identity

# Minikube support
make minikube-start  # Start minikube
make service-url     # Get service URL
```

## File Structure

```
├── app/auth-api/           # HTTP API layer
│   ├── handlers/           # HTTP request handlers
│   └── main.go            # Application entry point
├── business/              # Business logic layer
│   ├── store/             # Data access layer
│   └── validation/        # Input validation
├── foundation/            # Shared utilities
│   ├── middleware/        # HTTP middleware
│   ├── keys.go           # RSA key generation
│   ├── request.go        # Request utilities
│   └── response.go       # Response utilities
├── migrations/            # Database migration files
├── zarf/                 # Kubernetes/Helm deployment
│   ├── identity/         # Helm chart
│   └── *.yaml           # K8s manifests
├── Dockerfile            # Container build
├── Makefile             # Build automation
└── go.mod              # Go module definition
```

## Key Implementation Details

### Authentication Flow
1. User registers via `/register` endpoint
2. Password is hashed using bcrypt
3. User authenticates via `/login` endpoint
4. JWT token generated with RSA private key
5. Protected endpoints validate JWT with RSA public key

### Database Connection
- MySQL connection with connection pooling
- Automatic ping validation on startup
- Graceful connection handling and cleanup
- Migration support with golang-migrate

### Security Features
- Password hashing with bcrypt
- JWT tokens with RSA-256 signing
- CORS middleware support
- Input validation and sanitization

## Common Development Tasks

### Adding New Endpoints
1. Define route constants in `handlers/endpoints.go`
2. Implement handler function in appropriate handler file
3. Add route registration in `loadRoutes()` function
4. Add middleware as needed (Auth, CORS)

### Database Changes
1. Create new migration files in `migrations/` directory
2. Follow naming convention: `timestamp_description.up.sql` and `timestamp_description.down.sql`
3. Migrations run automatically on application startup

### Testing Guidelines
- Use testify for assertions
- Mock database connections with go-sqlmock
- Separate test environment configuration
- Test coverage currently at ~33%

## Build Commands

```bash
# Development
go run app/auth-api/main.go

# Testing
go test ./...

# Build
go build -o server app/auth-api/main.go

# Docker
docker build -t riyadennis/identity-server:1.3.0 .
```

## Monitoring & Health Checks
- **Liveness**: `GET /liveness` - Always returns 200 OK
- **Readiness**: `GET /readiness` - Checks database connectivity
- Logging with logrus (structured logging)
- Environment-based log levels

## CI/CD
- CircleCI integration (badge in README)
- Automated testing and coverage reporting
- Docker image publishing to Docker Hub
- Helm chart publishing to GitHub Pages

## Notes for Future Development
- Consider upgrading older dependencies (some from 2021-2022)
- Improve test coverage (currently 28-33%)
- Add input validation middleware
- Consider adding rate limiting
- Add comprehensive API documentation (OpenAPI/Swagger)
- Implement proper secret management for production
- Add database connection pooling configuration
- Consider adding metrics and observability