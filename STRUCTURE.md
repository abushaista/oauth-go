# OAuth Server - Complete Structure Summary

## ✅ Project Fully Implemented

This is a **production-ready OAuth 2.1 & OpenID Connect (OIDC) Server** built with Go following **Clean Architecture** and **CQRS** principles.

---

## 📁 Directory Structure

```
oauth-go/
│
├── cmd/api/
│   └── main.go                          # Entry point and Dependency Injection bootstrap
│
├── internal/
│   ├── domain/                          # Core Business Logic (Domain Layer)
│   │   ├── audit.go                     # Audit record structure
│   │   ├── audit_repository.go          # Audit persistence interface
│   │   ├── auth_code.go                 # AuthorizationCode model & repository interface
│   │   ├── client.go                    # Client application model & repository interface
│   │   ├── consent.go                   # Consent record model & repository interface
│   │   ├── crypto.go                    # Helper crypto constraints
│   │   ├── pkce.go                      # PKCE S256/Plain verification logic
│   │   ├── pkce_test.go                 # PKCE validator unit tests
│   │   ├── token.go                     # Access/Refresh token models & repositories interfaces
│   │   └── user.go                      # User model, lock logic, and repository interface
│   │
│   ├── application/                     # Application Business Logic (Application Layer)
│   │   ├── command/                     # Write CQRS Handlers (Mutations)
│   │   │   ├── interfaces.go            # Command types and command handler interfaces
│   │   │   ├── authorize_handler.go     # Creates short-lived authorization codes
│   │   │   ├── consent_handler.go       # Records user scope permissions for clients
│   │   │   ├── login_handler.go         # Authenticates users and handles lockout state
│   │   │   ├── refresh_handler.go       # Rotates refresh tokens and issues access tokens
│   │   │   ├── register_client_handler.go # Handles dynamic client registration (RFC 7591)
│   │   │   ├── revoke_handler.go        # Revokes access/refresh tokens (RFC 7009)
│   │   │   ├── token_handler.go         # Standard token exchanges (code, client_credentials)
│   │   │   └── utils.go                 # Safe secure random token generator
│   │   └── query/                       # Read CQRS Handlers (Data Retrieval)
│   │       ├── interfaces.go            # Query definitions and interfaces
│   │       ├── jwks_query.go            # Fetches public RSA keys
│   │       ├── userinfo_handler.go      # Prepares user profile payload for OIDC Userinfo
│   │       ├── client_handler.go        # Resolves client records
│   │       └── audit_query.go           # Fetches security events audit logs
│   │
│   ├── infrastructure/                  # Technical Implementations (Infrastructure Layer)
│   │   ├── persistence/                 # Persistence handlers
│   │   │   ├── audit_async.go           # Asynchronous worker to save audit logs safely
│   │   │   └── postgres/                # PostgreSQL repository implementations
│   │   │       ├── audit_repo.go
│   │   │       ├── auth_code_repo.go
│   │   │       ├── client_repo.go
│   │   │       ├── consent_repo.go
│   │   │       ├── refresh_token_repo.go
│   │   │       ├── token_repo.go
│   │   │       └── user_repo.go
│   │   └── security/                    # Security components
│   │       ├── jwks.go                  # Handles JSON Web Key Set keys and operations
│   │       ├── jwks_test.go
│   │       ├── jwt_rs256.go             # Creates and verifies RS256 JWT tokens
│   │       ├── jwt_rs256_test.go
│   │       ├── password.go              # Generates and validates bcrypt hashes
│   │       ├── password_test.go
│   │       └── refresh_hash.go          # Token utilities
│   │
│   └── interfaces/http/                 # HTTP API & Web Interface Layer
│       ├── handler_admin.go             # Client management API (admin protected)
│       ├── handler_audit.go             # Audits query endpoint (admin protected)
│       ├── handler_authorize.go         # GET /oauth/authorize endpoint (loads login/consent UI)
│       ├── handler_login.go             # POST /login (authenticates session cookie)
│       ├── handler_consent.go           # GET/POST /consent (renders UI / records scopes)
│       ├── handler_token.go             # POST /oauth/token endpoint
│       ├── handler_jwks.go              # GET /.well-known/jwks.json endpoint
│       ├── handler_oidc.go              # GET /.well-known/openid-configuration endpoint
│       ├── handler_registration.go      # POST /register (dynamic client registration)
│       ├── handler_revoke.go            # POST /oauth/revoke (token revocation)
│       ├── handler_static.go            # Serves local frontend files
│       ├── handler_userinfo.go          # GET /userinfo endpoint
│       ├── middleware_audit.go          # Audit log writer interceptor
│       ├── middleware_auth.go           # Bearer token validation middleware
│       ├── middleware_cors.go           # Cross-Origin resource sharing handler
│       ├── middleware_logging.go        # Structured HTTP request logger
│       ├── middleware_ratelimit.go      # IP-based sliding window rate limiter
│       ├── middleware_role.go           # User role-enforcement middleware
│       └── session.go                   # Signed HMAC-SHA256 cookie session manager
│
├── migrations/
│   └── init.sql                         # Database schema and seed data
│
├── tests/
│   └── integration_test.go              # Integration tests for core endpoints
│
├── web/                                 # User login and consent pages UI
│   ├── login.html
│   ├── consent.html
│   ├── style.css
│   ├── consent.css
│   ├── app.js
│   └── consent.js
│
├── go.mod                               # Go modules
├── go.sum                               # Checksum files
├── .gitignore                           # Git ignore configurations
├── .env.example                         # Environment configuration template
├── docker-compose.yml                   # Postgres, Server, and Adminer docker configurations
├── Makefile                             # Build & run tooling script
├── README.md                            # Main project readme
├── INDEX.md                             # Complete project documentation index
├── ARCHITECTURE.md                      # Architecture diagram and detail overview
├── COMPLETION.md                        # Completion log summary
└── QUICKSTART.md                        # Setup and test commands
```

---

## 🏗️ Architecture Layers

### Domain Layer (`internal/domain/`)
- **Entities**: User, Client, AuthorizationCode, AccessToken, RefreshToken, Consent, Audit
- **Repository Interfaces**: Abstract representations for databases (Postgres)
- **Business Rules**: PKCE verifier, user password validation, lockout rules, JWT creation.

### Application Layer (`internal/application/`)
- **CQRS Pattern**: Decoupled write actions (Commands) from read actions (Queries)
- **Command Handlers**: AuthorizeHandler, LoginHandler, ConsentHandler, RefreshHandler, RevokeHandler, RegisterClientHandler, TokenHandler
- **Query Handlers**: JWKSQuery, UserinfoHandler, ClientQueryHandler, AuditQueryHandler
- **No Direct tech integrations**: Relies strictly on Domain repository abstractions.

### Infrastructure Layer (`internal/infrastructure/`)
- **PostgreSQL Repositories**: Implementation of domain repositories via prepared SQL queries
- **Security Handlers**: BCrypt password hashing (Cost 14), RS256 token signer (custom JSON encoding), JWKS keys generator
- **Asynchronous Writer**: Asynchronous audit logging helper that accepts logs and writes them inside a concurrent background queue.

### Interface Layer (`internal/interfaces/http/`)
- **Controllers**: Adapts HTTP requests to application commands and queries
- **Middlewares**: Handles request rate limiting, CORS configuration, bearer auth extraction, logging and request auditing.
- **Sessions**: Uses cookie storage signed with an HMAC-SHA256 signature to manage state across the login/consent redirection.

---

## 🌐 API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/health` | GET | Public | Server health checks |
| `/oauth/authorize` | GET/POST | Session | Authorization endpoint (loads Login/Consent UI) |
| `/oauth/token` | POST | Client Secret | Token exchange (auth code, refresh, client credentials) |
| `/oauth/revoke` | POST | Client Secret | Revokes access and refresh tokens (RFC 7009) |
| `/.well-known/jwks.json` | GET | Public | Standard JSON Web Key Set discovery |
| `/.well-known/openid-configuration` | GET | Public | OpenID Connect metadata discovery |
| `/login` | POST | Public | Authenticates user credentials & sets session cookie |
| `/consent` | GET/POST | Session | Renders scope options / registers consent |
| `/userinfo` | GET | Bearer Token | Returns OIDC profile details (sub, role, username) |
| `/register` | POST | Public | Standard dynamic client registration (RFC 7591) |
| `/audits` | GET | Admin Bearer | Fetches system logs matching user or client query |
| `/admin/api` | GET | Admin Bearer | Client administration API |

---

## 🗄️ Database Schema

Schema configured with strict foreign keys and index trees:
- **`users`**: user logins, role (`admin` or `user`), `failed_login_attempts`, and account lock `locked_until`
- **`clients`**: registration entries, client IDs, client secrets, and redirect URIs
- **`authorization_codes`**: short-lived codes with scopes and PKCE fields
- **`access_tokens`**: bearer tokens mapping user and client
- **`refresh_tokens`**: refresh tokens with tracking revoked state (`revoked` true/false)
- **`consents`**: user-approved client applications and scopes
- **`audit_logs`**: structured record of security and authorization events

---

## 📦 Dependencies

Managed by Go modules:
- `github.com/lib/pq`: PostgreSQL driver
- `golang.org/x/crypto/bcrypt`: Cryptographic hashing for passwords

---

## ✨ Features Implemented

✅ OAuth 2.0 / 2.1 Authorization Code Flow with enforced PKCE  
✅ Client Credentials Grant Flow  
✅ OpenID Connect (OIDC) metadata discovery & `/userinfo` endpoint  
✅ Custom JWT RS256 signing and dynamic JWKS key set  
✅ Refresh Token Rotation (RTR) and revocation endpoints  
✅ Safe BCrypt password hashing & account lockout protection (5 fails)  
✅ Non-blocking asynchronous audit logging channel  
✅ Global CORS, request logs, and IP-based rate limiting middlewares  
✅ UI templates built in Vue.js (served locally via standard Go files)  

---

## 📊 Database Relationships

```
users (1) ──→ (many) authorization_codes
users (1) ──→ (many) access_tokens
users (1) ──→ (many) refresh_tokens
users (1) ──→ (many) consents
users (1) ──→ (many) audit_logs

clients (1) ──→ (many) authorization_codes
clients (1) ──→ (many) access_tokens
clients (1) ──→ (many) refresh_tokens
clients (1) ──→ (many) consents
clients (1) ──→ (many) audit_logs
```

---

## ✅ Completion Status

**Status**: 🟢 **FULLY IMPLEMENTED & PRODUCTION-READY**

Every Clean Architecture and CQRS layer is complete and verified under standard Go 1.26.1 and Docker environments.
