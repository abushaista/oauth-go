# OAuth Server

A production-ready OAuth 2.0 authorization server built with Go, following Clean Architecture principles.

## Features

- **OAuth 2.0 Authorization Code Flow** with PKCE support
- **JWT Token Generation** with RS256 algorithm
- **PostgreSQL** persistence layer
- **Clean Architecture** with domain, application, infrastructure, and interface layers
- **Audit Logging** for security and compliance
- **User Consent Management** for OAuth authorization
- **Refresh Token Support** for long-lived sessions

## Project Structure

```
oauth-server/
├── cmd/api/
│   └── main.go              # Entry point
├── internal/
│   ├── domain/              # Domain entities and interfaces
│   │   ├── user.go
│   │   ├── client.go
│   │   ├── auth_code.go
│   │   ├── token.go
│   │   ├── consent.go
│   │   ├── audit.go
│   │   ├── pkce.go
│   │   ├── crypto.go
│   │   └── *_repository.go  # Repository interfaces
│   ├── application/         # Application business logic
│   │   ├── command/         # Command handlers (CQRS pattern)
│   │   │   ├── authorize_handler.go
│   │   │   ├── token_handler.go
│   │   │   ├── refresh_handler.go
│   │   │   └── login_handler.go
│   │   └── query/           # Query handlers
│   │       └── jwks_query.go
│   ├── infrastructure/      # External service implementations
│   │   ├── persistence/postgres/
│   │   │   ├── user_repo.go
│   │   │   ├── client_repo.go
│   │   │   ├── auth_code_repo.go
│   │   │   ├── token_repo.go
│   │   │   ├── refresh_token_repo.go
│   │   │   ├── consent_repo.go
│   │   │   └── audit_repo.go
│   │   └── security/
│   │       ├── password.go
│   │       ├── jwt_rs256.go
│   │       ├── jwks.go
│   │       └── refresh_hash.go
│   └── interfaces/http/     # HTTP handlers
│       ├── handler_authorize.go
│       ├── handler_login.go
│       ├── handler_consent.go
│       ├── handler_token.go
│       ├── handler_jwks.go
│       └── handler_static.go
├── migrations/
│   └── init.sql             # Database schema
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/abushaista/oauth-go.git
cd oauth-go
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
# Create database
createdb oauth_db

# Run migrations
psql oauth_db < migrations/init.sql
```

4. Set environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/oauth_db?sslmode=disable"
export PORT=8080
```

5. Run the server:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### OAuth Endpoints

- `GET /oauth/authorize` - Authorization code request
- `POST /oauth/token` - Token exchange endpoint
- `GET /.well-known/jwks.json` - Public key set

### User Endpoints

- `POST /login` - User login
- `GET/POST /consent` - User consent management

### Health Check

- `GET /health` - Server health check

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://user:password@localhost:5432/oauth_db?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | HTTP server port |

## Database Schema

### Users
- `id` - Unique user identifier
- `username` - Username for login
- `password` - Hashed password

### Clients
- `id` - Unique client identifier
- `client_id` - OAuth client ID
- `client_secret` - OAuth client secret
- `redirect_uri` - Authorized redirect URI

### Authorization Codes
- `code` - Authorization code
- `user_id` - Associated user
- `client_id` - Associated client
- `code_challenge` - PKCE code challenge
- `code_challenge_method` - PKCE method (plain or S256)
- `expires_at` - Expiration time

### Tokens
- **Access Tokens**: Short-lived tokens for API access
- **Refresh Tokens**: Long-lived tokens for getting new access tokens

### Consents
- Stores user consent for client applications
- Tracks which scopes are authorized

### Audit Logs
- Logs all significant authentication events
- Includes user IP addresses and action details

## Development

### Running Tests

```bash
go test ./...
```

### Code Style

Follow standard Go conventions:
```bash
go fmt ./...
go vet ./...
```

## Architecture Patterns

This project implements several architecture patterns:

### Clean Architecture
- **Domain Layer**: Business rules and entities
- **Application Layer**: Use cases and business logic
- **Infrastructure Layer**: Technical implementations
- **Interface Layer**: External interfaces (HTTP, gRPC, etc.)

### CQRS Pattern
- Commands for state-changing operations (login, authorize, token exchange)
- Queries for read operations (JWKS retrieval)

### Repository Pattern
- Abstract data access behind interfaces
- Easy to swap implementations (PostgreSQL, MongoDB, Redis, etc.)

## Security Considerations

- Password hashing with bcrypt (implement in `password.go`)
- JWT signing with RS256 asymmetric algorithm
- PKCE support for public clients
- Secure refresh token handling
- Audit logging for compliance

## TODO

- [ ] Implement password hashing with bcrypt
- [ ] Complete JWT signing with RS256
- [ ] Add rate limiting
- [ ] Add CORS support
- [ ] Add OpenID Connect support
- [ ] Add client refresh token rotation
- [ ] Add token revocation endpoints
- [ ] Add user registration endpoint
- [ ] Add React frontend for login/consent screens

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
