# identity-server
![Coverage](https://img.shields.io/badge/Coverage-63.8%25-yellow)

## Project Overview
A Go-based identity management server providing user registration, authentication, and JWT token generation. Built with a clean architecture separating business logic, storage, and HTTP handling layers.

## Architecture

## Core Components
- **Main Entry Point**: `app/auth-api/main.go` - Application startup, database connection, and server initialization
- **HTTP Handlers**: `app/auth-api/rest/` - REST API endpoints and routing
- **Server Layer**: `app/auth-api/server/` - HTTP server configuration and lifecycle management
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
- **HTTP Router**: go-chi/chi/v5
- **Migrations**: golang-migrate/migrate
- **Documentation**: Swagger/OpenAPI (swaggo/swag)
- **Deployment**: Docker + Kubernetes/Helm

## API Endpoints

### Public Endpoints
- `POST /register` - User registration
- `POST /login` - User authentication and token generation
- `GET /liveness` - Kubernetes liveness probe
- `GET /readiness` - Kubernetes readiness probe

### Protected Endpoints (require JWT)
- `GET /user/home` - User profile access
- `DELETE /admin/delete/:id` - User deletion

To Run the service locally we need .env file set with the following values:

```bash
PORT="8089"
ENV="dev"
MYSQL_DATABASE="identity"
MYSQL_USERNAME="username"
MYSQL_PASSWORD="password"
MYSQL_PORT="3306"
MYSQL_HOST="127.0.0.1"
MIGRATION_PATH="migrations"
```