# OAuth Server - Complete File Index

## 📖 Documentation Files (Start Here!)

| File | Purpose |
|------|---------|
| **[README.md](README.md)** | 📘 Full project documentation |
| **[QUICKSTART.md](QUICKSTART.md)** | ⚡ 5-minute setup guide |
| **[STRUCTURE.md](STRUCTURE.md)** | 🏗️ Project structure overview |
| **[ARCHITECTURE.md](ARCHITECTURE.md)** | 🎨 Visual architecture diagrams |
| **[COMPLETION.md](COMPLETION.md)** | ✅ Project completion summary |
| **INDEX.md** | 📑 This file - complete index |

---

## 🚀 Getting Started (Pick One)

### Option 1: Fastest (Recommended)
```bash
make dev           # Start postgres db in docker + run Go server locally
# Server on http://localhost:8080
# Database UI on http://localhost:8081
```

### Option 2: Pure Docker setup
```bash
make docker-up     # Starts Postgres + Server + Adminer in containers
```

### Option 3: Manual
See **[QUICKSTART.md](QUICKSTART.md)** for detailed instructions.

---

## 📁 Project File Structure

### Configuration
```
.env.example          Configuration template (copy to .env)
.gitignore           Git ignore patterns
docker-compose.yml   PostgreSQL + Server + Adminer setup
Dockerfile           Multi-stage Docker build configuration
go.mod               Go module definition
go.sum               Go checksum values
Makefile             Build commands
```

### Documentation
```
README.md            Main documentation
QUICKSTART.md        Quick start guide
STRUCTURE.md         Project structure details
ARCHITECTURE.md      Architecture diagrams
COMPLETION.md        Completion summary
INDEX.md             This file
```

### Application Code

#### Entry Point
```
cmd/api/main.go      Application entry point & dependency injection
```

#### Domain Layer (Core Business Logic)
```
internal/domain/
├── user.go                    User entity + UserRepository interface
├── client.go                  Client entity + ClientRepository interface
├── auth_code.go               AuthorizationCode entity + interface
├── token.go                   AccessToken & RefreshToken entities + repositories
├── consent.go                 Consent entity + interface
├── audit.go                   Audit entity
├── audit_repository.go        Audit repository interface
├── pkce.go                    PKCE validation logic
└── crypto.go                  Crypto constants
```

#### Application Layer (CQRS Pattern)
```
internal/application/
├── command/
│   ├── interfaces.go          Command interfaces and definitions
│   ├── authorize_handler.go   Authorization code flow handler
│   ├── token_handler.go       Token exchange handler
│   ├── refresh_handler.go     Refresh token handler
│   ├── login_handler.go       User login and lockout handler
│   ├── consent_handler.go     User consent command handler
│   ├── revoke_handler.go      Token revocation handler
│   ├── register_client_handler.go Dynamic client registration handler
│   └── utils.go               Cryptographic helpers
└── query/
    ├── interfaces.go          Query interfaces and definitions
    ├── jwks_query.go          JWKS retrieval handler
    ├── userinfo_handler.go    Userinfo retrieval handler
    ├── client_handler.go      Client lookup query handler
    └── audit_query.go         Audit logs retrieval query handler
```

#### Infrastructure Layer (Database & Security)
```
internal/infrastructure/
├── persistence/
│   ├── audit_async.go         Asynchronous audit worker channel
│   └── postgres/
│       ├── user_repo.go           PostgreSQL user repository
│       ├── client_repo.go         PostgreSQL client repository
│       ├── auth_code_repo.go      PostgreSQL auth code repository
│       ├── token_repo.go          PostgreSQL access token repository
│       ├── refresh_token_repo.go  PostgreSQL refresh token repository
│       ├── consent_repo.go        PostgreSQL consent repository
│       └── audit_repo.go          PostgreSQL audit repository
└── security/
    ├── password.go            Password hashing (bcrypt cost 14)
    ├── jwt_rs256.go           JWT signing with RS256 (RSA keys)
    ├── jwks.go                JSON Web Key Set provider
    └── refresh_hash.go        PKCE and token utilities
```

#### Interfaces Layer (HTTP API)
```
internal/interfaces/http/
├── handler_admin.go          Client management API (admin only)
├── handler_audit.go          Audit query API
├── handler_authorize.go      GET/POST /oauth/authorize
├── handler_login.go          POST /login
├── handler_consent.go        GET/POST /consent
├── handler_token.go          POST /oauth/token
├── handler_jwks.go           GET /.well-known/jwks.json
├── handler_oidc.go           GET /.well-known/openid-configuration
├── handler_registration.go   POST /register
├── handler_revoke.go         POST /oauth/revoke
├── handler_static.go         Static UI server
├── handler_userinfo.go       GET /userinfo
├── middleware_audit.go       Audit logger middleware
├── middleware_auth.go        Bearer token authentication
├── middleware_cors.go        CORS middleware
├── middleware_logging.go     Request logging middleware
├── middleware_ratelimit.go   IP-based rate limiter middleware
├── middleware_role.go        User role validation middleware
└── session.go                Secure HMAC-signed sessions
```

#### Database
```
migrations/
└── init.sql                   Complete database schema + seeds
```

#### Frontend
```
web/
├── login.html                 Login view
├── consent.html               Consent view
├── style.css                  Styles for login
├── consent.css                Styles for consent
├── app.js                     Vue.js login scripts
└── consent.js                 Vue.js consent scripts
```

---

## 🔍 Quick File Reference

### By Feature

#### Authentication & Account Lockout
- `domain/user.go` - User model and repository
- `infrastructure/persistence/postgres/user_repo.go` - User DB persistence
- `infrastructure/security/password.go` - Password handling (bcrypt cost 14)
- `application/command/login_handler.go` - Login & account lockout business logic
- `interfaces/http/handler_login.go` - Login API controller
- `interfaces/http/session.go` - Session manager (signed cookie sessions)

#### Authorization & PKCE
- `domain/auth_code.go` - Authorization code model
- `domain/pkce.go` - PKCE validation logic
- `infrastructure/persistence/postgres/auth_code_repo.go` - Auth code DB persistence
- `application/command/authorize_handler.go` - Code authorization logic
- `interfaces/http/handler_authorize.go` - Authorization API controller

#### Token Management & Rotation
- `domain/token.go` - Token models and repository interfaces
- `infrastructure/persistence/postgres/token_repo.go` - Access token DB persistence
- `infrastructure/persistence/postgres/refresh_token_repo.go` - Refresh token DB persistence
- `infrastructure/security/jwt_rs256.go` - JWT RS256 token encoding and signing
- `application/command/token_handler.go` - Token exchange logic
- `application/command/refresh_handler.go` - Token refresh and rotation logic
- `application/command/revoke_handler.go` - Token revocation logic
- `interfaces/http/handler_token.go` - Token endpoint controller
- `interfaces/http/handler_revoke.go` - Token revocation controller

#### OpenID Connect (OIDC)
- `interfaces/http/handler_oidc.go` - OIDC discovery provider metadata endpoint
- `interfaces/http/handler_userinfo.go` - User details provider endpoint
- `application/query/userinfo_handler.go` - User profile query lookup

#### Consent Management
- `domain/consent.go` - Consent model
- `infrastructure/persistence/postgres/consent_repo.go` - Consent DB persistence
- `interfaces/http/handler_consent.go` - Consent API controller

#### Audit & Logging
- `domain/audit.go` - Audit model and interface
- `infrastructure/persistence/audit_async.go` - Thread-safe asynchronous logging worker
- `infrastructure/persistence/postgres/audit_repo.go` - Audit DB persistence
- `interfaces/http/middleware_audit.go` - Route action audit logger middleware
- `interfaces/http/handler_audit.go` - Query audit history endpoint

#### Key Management
- `infrastructure/security/jwks.go` - JWKS private/public key set generator
- `interfaces/http/handler_jwks.go` - Public keys discovery endpoint

---

## 📊 Statistics

| Metric | Count |
|--------|-------|
| Go source files | 45+ |
| Domain models | 9 |
| Repository interfaces | 8 |
| Application handlers | 12 |
| HTTP handlers | 12 |
| Database tables | 7 |
| API endpoints | 10+ |
| Lines of code | 3000+ |
| Documentation files | 6 |

---

## 🎯 API Endpoints

```
GET  /health                            Health check
GET  /oauth/authorize                   Authorization request (triggers login/consent UI)
POST /oauth/token                       Token exchange / client credentials / refresh
POST /oauth/revoke                      Token revocation
GET  /.well-known/jwks.json             Public key set
GET  /.well-known/openid-configuration  OIDC discovery configuration
POST /login                             User session login
GET  /consent                           User consent UI loader
POST /consent                           Consent scopes submission
GET  /userinfo                          User details (requires bearer token)
POST /register                          Dynamic client registration
GET  /audits                            Security audit logs (Admin only)
GET  /admin/api                         Admin client query API (Admin only)
```

---

## 🗄️ Database Tables

```
users                   User credentials + failed logins + lock times
clients                 OAuth client definitions
authorization_codes    Short-lived authorization codes & PKCE params
access_tokens          Issued active access tokens
refresh_tokens         Issued active refresh tokens
consents               User approved client scopes
audit_logs             Structured security logs
```

---

## 🛠️ Available Commands

```bash
make help              Show all commands
make build             Build server binary
make run               Build and run server locally
make test              Run unit & integration tests
make fmt               Format Go code
make lint              Run golangci-lint analysis
make clean             Clean build artifacts
make docker-up         Start PostgreSQL, Server, and Adminer in Docker
make docker-down       Stop and remove docker containers
make docker-db         Start only Postgres + Adminer in Docker
make dev               Start local Go server connected to Docker Postgres database
```

---

## 📚 Reading Order

1. **Start here**: [README.md](README.md)
2. **Setup**: [QUICKSTART.md](QUICKSTART.md)
3. **Understanding**: [STRUCTURE.md](STRUCTURE.md)
4. **Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md)
5. **Completion**: [COMPLETION.md](COMPLETION.md)

---

## 🔑 Key Technologies

- **Language**: Go 1.26.1
- **Libraries**: standard `net/http` package (no heavy framework overhead)
- **Database**: PostgreSQL 15+
- **Driver**: `github.com/lib/pq`
- **Crypto**: `golang.org/x/crypto/bcrypt`
- **UI Framework**: Vue.js (for Login/Consent page scripts)
- **Pattern**: Clean Architecture + CQRS + Repository Pattern

---

## ✨ Features Implemented

✅ OAuth 2.0 / 2.1 Authorization Code Flow  
✅ Enforced PKCE Support  
✅ JWT Tokens Signed with RS256  
✅ JSON Web Key Set (JWKS) Publication  
✅ OpenID Connect (OIDC) Provider metadata  
✅ Userinfo Profile Endpoint  
✅ Refresh Token Rotation & Revocation  
✅ User Authentication & HMAC-signed Session management  
✅ Account Lockout Policy (5 failed attempts lock for 15 minutes)  
✅ Dynamic Client Registration (RFC 7591)  
✅ Asynchronous Audit Logging Worker  
✅ IP-based Global Rate Limiting  
✅ Global CORS Middleware  
✅ Postgres persistence with optimized pooling  
✅ Vue.js login/consent interfaces  

---

## 🚀 Next Steps / Enhancement Ideas

- [ ] Add support for Authorization Code Flow with hybrid response types.
- [ ] Implement database-backed user registration endpoint.
- [ ] Add support for mutual TLS client authentication (mTLS).
- [ ] Integrate Prometheus telemetry metrics endpoint.
- [ ] Setup a full CI/CD pipeline script.

---

**Last Updated**: May 27, 2026  
**Status**: ✅ Complete and Production-ready  
**Version**: 2.0.0  
