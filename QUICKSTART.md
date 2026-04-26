# Quick Start Guide

## Setup Development Environment

### Option 1: Using Docker (Recommended)

1. Start PostgreSQL and Adminer:
```bash
make docker-up
```

2. Run database migrations:
```bash
make migrate
```

3. Start the development server:
```bash
make dev
```

The server will be available at `http://localhost:8080`
Adminer (database UI) will be available at `http://localhost:8081`

### Option 2: Manual Setup

1. Create PostgreSQL database:
```bash
createdb oauth_db
```

2. Create user:
```bash
createuser -P oauth_user
# Enter password when prompted
```

3. Grant privileges:
```bash
psql -d oauth_db -c "GRANT ALL PRIVILEGES ON DATABASE oauth_db TO oauth_user;"
```

4. Run migrations:
```bash
psql -h localhost -U oauth_user -d oauth_db -f migrations/init.sql
```

5. Set environment variables:
```bash
export DATABASE_URL="postgres://oauth_user:YOUR_PASSWORD@localhost:5432/oauth_db?sslmode=disable"
export PORT=8080
```

6. Run the server:
```bash
go run cmd/api/main.go
```

## Available Commands

```bash
make build      # Build binary
make run        # Run the server
make test       # Run tests
make fmt        # Format code
make lint       # Run linter
make clean      # Clean artifacts
make docker-up  # Start Docker containers
make docker-down # Stop Docker containers
make migrate    # Run migrations
make dev        # Start development environment
```

## API Testing

### Health Check
```bash
curl http://localhost:8080/health
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass"
  }'
```

### Authorization Request
```bash
curl "http://localhost:8080/oauth/authorize?client_id=client1&redirect_uri=http://localhost:3000/callback&response_type=code&scope=openid+profile&state=state123"
```

### Token Exchange
```bash
curl -X POST http://localhost:8080/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code&code=AUTH_CODE&client_id=client1&client_secret=secret1&redirect_uri=http://localhost:3000/callback"
```

### Get JWKS
```bash
curl http://localhost:8080/.well-known/jwks.json
```

## Database Management

### Using Adminer
- URL: `http://localhost:8081`
- Server: `postgres`
- Username: `oauth_user`
- Password: `oauth_password`
- Database: `oauth_db`

### Using psql CLI
```bash
# Connect to database
psql -h localhost -U oauth_user -d oauth_db

# List tables
\dt

# Describe a table
\d users

# Exit
\q
```

## Next Steps

1. **Password Hashing**: Implement bcrypt in `internal/infrastructure/security/password.go`
2. **JWT Signing**: Implement RS256 in `internal/infrastructure/security/jwt_rs256.go`
3. **Unit Tests**: Add tests for handlers and repositories
4. **API Documentation**: Add Swagger/OpenAPI documentation
5. **React Frontend**: Create login and consent screens in `web/` directory
6. **Rate Limiting**: Implement rate limiting middleware
7. **CORS Support**: Add CORS middleware
8. **OpenID Connect**: Implement OIDC features

## Troubleshooting

### Database Connection Error
```
Error: could not connect to server: No such file or directory
```
Make sure PostgreSQL is running and `DATABASE_URL` is set correctly.

### Port Already in Use
Change the port in `.env` file or run:
```bash
PORT=8081 go run cmd/api/main.go
```

### Migration Errors
Drop and recreate the database:
```bash
dropdb oauth_db
createdb oauth_db
make migrate
```

## Project Structure Overview

- **cmd/api/**: Application entry point
- **internal/domain/**: Core business logic and domain models
- **internal/application/**: Use cases and handlers (CQRS pattern)
- **internal/infrastructure/**: Database, security, external services
- **internal/interfaces/http/**: HTTP API handlers
- **migrations/**: Database schema
