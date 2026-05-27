# OAuth 2.1 & OpenID Connect (OIDC) Server

A production-ready, strict OAuth 2.1 and OpenID Connect (OIDC) authorization server built with Go, following Clean Architecture and CQRS principles.

## 🚀 Features

- **Strict OAuth 2.1 Compliance**: Enforces PKCE (Proof Key for Code Exchange) for public clients, and disables insecure legacy flows (Implicit and Resource Owner Password Credentials).
- **OpenID Connect (OIDC) Support**: Includes OIDC Discovery (`/.well-known/openid-configuration`), ID Token generation signed with RS256, and a `/userinfo` endpoint.
- **Asymmetric Cryptography (RS256 & JWKS)**: Cryptographically signs JWTs using RS256 with key rotated sets published via the JSON Web Key Set (JWKS) endpoint.
- **Dynamic Client Registration**: Standardized endpoint (`POST /register`) for registering clients on the fly (RFC 7591).
- **Token Revocation**: Standardized endpoint (`POST /oauth/revoke`) for revoking access and refresh tokens (RFC 7009).
- **Refresh Token Rotation**: Automatically rotates refresh tokens on every usage, invalidating the old token to mitigate replay attacks.
- **Robust Account Lockout Policy**: Tracks failed logins; locks a user account for 15 minutes after 5 consecutive failed attempts.
- **High-Performance Asynchronous Auditing**: Logs all security events asynchronously using a worker channel pool to avoid blocking HTTP response flows.
- **Global Middlewares**: Enforces CORS configuration, IP-based rate limiting (10,000 requests/minute), and structured logging.
- **Clean Architecture & CQRS**: Strictly decoupled layers separating Domain entities, CQRS Commands/Queries, Infrastructure, and HTTP Interfaces.
- **Docker-Ready**: Packaged with multi-stage Docker builds and automated Postgres initialization.

## 📁 Project Structure

```
oauth-go/
├── cmd/api/
│   └── main.go                         # Entry point & dependency injection
├── internal/
│   ├── domain/                         # Domain Layer (Business entities & interfaces)
│   │   ├── audit.go                    # Audit entity
│   │   ├── audit_repository.go         # Audit Repository interface
│   │   ├── auth_code.go                # Auth Code entity & interface
│   │   ├── client.go                   # Client entity & interface
│   │   ├── consent.go                  # Consent entity & interface
│   │   ├── crypto.go                   # Crypto helper constants
│   │   ├── pkce.go                     # PKCE validation logic
│   │   ├── pkce_test.go                # PKCE unit tests
│   │   ├── token.go                    # Access/Refresh Token entities & repositories
│   │   └── user.go                     # User entity & UserRepository interface
│   ├── application/                    # Application Layer (CQRS Use cases)
│   │   ├── command/                    # Write commands (State-mutators)
│   │   │   ├── authorize_handler.go    # Issue authorization code
│   │   │   ├── consent_handler.go      # Submit user consents
│   │   │   ├── interfaces.go           # Command types & definitions
│   │   │   ├── login_handler.go        # Handle user login and lockout
│   │   │   ├── refresh_handler.go      # Rotate and refresh tokens
│   │   │   ├── register_client_handler.go # Register dynamic clients
│   │   │   ├── revoke_handler.go       # Revoke access/refresh tokens
│   │   │   ├── token_handler.go        # Exchange code/credentials for tokens
│   │   │   └── utils.go                # Cryptographic helper utilities
│   │   └── query/                      # Read queries (Retrievals)
│   │       ├── audit_query.go          # Retrieve audit history
│   │       ├── client_handler.go       # Fetch client details
│   │       ├── interfaces.go           # Query types & definitions
│   │       ├── jwks_query.go           # Retrieve JWKS
│   │       └── userinfo_handler.go     # Fetch OIDC user info
│   ├── infrastructure/                 # Infrastructure Layer (Tech implementations)
│   │   ├── persistence/                # Persistence implementations
│   │   │   ├── audit_async.go          # Non-blocking async writer
│   │   │   └── postgres/               # PostgreSQL concrete repositories
│   │   │       ├── audit_repo.go
│   │   │       ├── auth_code_repo.go
│   │   │       ├── client_repo.go
│   │   │       ├── consent_repo.go
│   │   │       ├── refresh_token_repo.go
│   │   │       ├── token_repo.go
│   │   │       └── user_repo.go
│   │   └── security/                   # Core security components
│   │       ├── jwks.go                 # JWKS dynamic key provider
│   │       ├── jwks_test.go
│   │       ├── jwt_rs256.go            # Custom RS256 token signer
│   │       ├── jwt_rs256_test.go
│   │       ├── password.go             # BCrypt password hasher
│   │       ├── password_test.go
│   │       └── refresh_hash.go         # Safe hashing helpers
│   └── interfaces/http/                # Interface Layer (API & UI controllers)
│       ├── handler_admin.go            # Protected admin controls
│       ├── handler_audit.go            # View client/user audit logs
│       ├── handler_authorize.go        # User authorization page endpoint
│       ├── handler_consent.go          # Handle client consent flow
│       ├── handler_jwks.go             # Expose public keys
│       ├── handler_login.go            # Authenticate user session
│       ├── handler_oidc.go             # OpenID provider configuration
│       ├── handler_registration.go     # Dynamic client registration API
│       ├── handler_revoke.go           # Token revocation API
│       ├── handler_static.go           # Serves HTML pages
│       ├── handler_token.go            # Token endpoints
│       ├── handler_userinfo.go         # Fetch user details using bearer tokens
│       ├── middleware_audit.go         # Automatically records action audits
│       ├── middleware_auth.go          # Decodes and validates bearer tokens
│       ├── middleware_cors.go          # Handles cross-origin requests
│       ├── middleware_logging.go       # Formats HTTP logs
│       ├── middleware_ratelimit.go     # Limits request rates
│       ├── middleware_role.go          # Limits handlers by user role (admin)
│       └── session.go                  # HMAC-signed session manager
├── migrations/
│   └── init.sql                        # PostgreSQL schema & seeding
├── tests/
│   └── integration_test.go             # HTTP API handler integration tests
├── web/                                # Vue.js Frontend Client UI
│   ├── login.html                      # Login View
│   ├── consent.html                    # Consent View
│   ├── style.css                       # Styles for Login
│   ├── consent.css                     # Styles for Consent
│   ├── app.js                          # Login Script
│   └── consent.js                      # Consent Script
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ Getting Started

### Prerequisites

- **Go 1.26.1**
- **Docker & Docker Compose** (for easiest configuration)

### Running with Docker (Recommended)

To start the database and the OAuth server, simply run:

```bash
make docker-up
```

- **OAuth Server**: runs at `http://localhost:8080`
- **Database UI (Adminer)**: runs at `http://localhost:8081`

### Local Development (Local Server + Dockerized DB)

If you prefer to run the Go application locally:

```bash
make dev
```
This command starts the PostgreSQL database inside Docker, waits for it to become ready, and then starts the Go server locally.

To stop the background containers, run:
```bash
make docker-down
```

## 🎯 API Endpoints

### OAuth 2.1 & OIDC Endpoints
- `GET /oauth/authorize` - Authorization Request (triggers login/consent pages)
- `POST /oauth/token` - Token exchange (code, client credentials, refresh token)
- `POST /oauth/revoke` - Revokes access/refresh tokens
- `GET /.well-known/jwks.json` - Returns the JWKS containing public keys
- `GET /.well-known/openid-configuration` - OpenID Connect Provider metadata
- `GET /userinfo` - Fetches authenticated user info using bearer token

### User Session & Management Endpoints
- `POST /login` - Direct user session login (verifies bcrypt hash)
- `GET/POST /consent` - Fetches/Submits scopes consent for a client application
- `POST /register` - Dynamic client registration (returns client secret & ID)

### Administration & Telemetry
- `GET /audits` - Fetch security audit trails (Admin only)
- `GET/POST /admin/api` - Client management controls (Admin only)
- `GET /health` - Service health status

## ⚙️ Environment Variables

Copy `.env.example` to `.env` or set these in your environment:

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://oauth_user:oauth_password@localhost:5432/oauth_db?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | HTTP Port for the OAuth server |
| `JWT_SECRET` | `dev-secret-key` | Key used for signing cookie sessions |

## 🗄️ Database Schema

The database consists of 7 tables configured in `migrations/init.sql`:

1. **`users`**: Stores user accounts, bcrypt password hashes, and lockout data (`failed_login_attempts`, `locked_until`).
2. **`clients`**: Stores client client ID, hashed secret, and authorized redirect URIs.
3. **`authorization_codes`**: Short-lived authorization codes, scopes, and PKCE challenge methods.
4. **`access_tokens`**: Active JWT access tokens.
5. **`refresh_tokens`**: Active refresh tokens tracked for revocation and rotation.
6. **`consents`**: Persists user authorizations of specific scopes for client applications.
7. **`audit_logs`**: Chronological track of security actions, IP addresses, and actors.

## 🧪 Development and Testing

### Running Tests

Run all unit and integration tests:
```bash
make test
```

### Formatter & Linter
```bash
make fmt
make lint
```

## 🔒 Security Considerations

1. **PKCE Enforcement**: Enforced for all authorization code flow requests. Insecure code flow requests without `code_challenge` are automatically rejected.
2. **Brute Force Lockout**: Accounts are locked for 15 minutes upon 5 failed logins. Attempts during lockout return a `403 Forbidden` immediately.
3. **Bcrypt Cost**: Hashed using BCrypt with a high work factor (Cost: 14).
4. **Refresh Rotation**: Ensures any refresh token is single-use only. Using a rotated token triggers a fraud event, logging the incident and invalidating associated active tokens.
5. **Asymmetric Signing**: Keys are kept in-memory and public verification keys are rotated and published via the standard JWKS structure.

## 📄 License

MIT
