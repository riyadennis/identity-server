# identity-server
![Coverage](https://img.shields.io/badge/Coverage-18.4%25-red)

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

#### Register
```bash
curl -X POST http://localhost:8089/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Doe",
    "email": "jane.doe@example.com",
    "password": "SecurePassword123!",
    "company": "Acme Corp",
    "post_code": "SW1A 1AA",
    "terms": true
  }'
```

#### Login
Login uses HTTP Basic Auth (email:password):
```bash
curl -X POST http://localhost:8089/login \
  -u "jane.doe@example.com:SecurePassword123!"
```

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
### Login Mutation example
```
mutation Login($input: LoginInput!) { Login(input: $input) { status accessToken expiry tokenType lastRefresh tokenTTL } }
```

## GraphQL API Documentation

The GraphQL API docs are generated using [SpectaQL](https://github.com/anvilco/spectaql) from the schema at `app/gql/graph/schema.graphqls`.

### Install SpectaQL

```bash
npm install -g spectaql
```

### Generate docs

```bash
npx spectaql app/gql/spectaql.yaml -t docs/graphql
```

This outputs static HTML documentation to the `docs/graphql/` directory.

### Preview docs locally

```bash
npx spectaql app/gql/spectaql.yaml -t docs/graphql -D
```

The `-D` flag starts a development server so you can view the docs in your browser.
