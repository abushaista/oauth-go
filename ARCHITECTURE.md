# OAuth Server Architecture Visualization

## Request Flow Diagram

```
HTTP Request
    │
    ▼
┌────────────────────────────────────────────────────────┐
│     HTTP Interface Layer                               │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Router (http.ServeMux)                           │  │
│  │ ├─ GET  /health                                  │  │
│  │ ├─ GET/POST /oauth/authorize                     │  │
│  │ ├─ POST /oauth/token                             │  │
│  │ ├─ POST /oauth/revoke                            │  │
│  │ ├─ GET  /.well-known/jwks.json                   │  │
│  │ ├─ GET  /.well-known/openid-configuration        │  │
│  │ ├─ POST /login                                   │  │
│  │ ├─ GET/POST /consent                             │  │
│  │ ├─ GET  /userinfo                                │  │
│  │ ├─ POST /register                                │  │
│  │ ├─ GET  /audits (Admin Only)                     │  │
│  │ └─ GET/POST /admin/api (Admin Only)              │  │
│  └──────────────────────────────────────────────────┘  │
│              │                                         │
│              ▼                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Middlewares & Filters                            │  │
│  │ ├─ requestLogger.Wrap()                          │  │
│  │ ├─ corsMiddleware.Wrap()                         │  │
│  │ ├─ rateLimiter.Wrap()                            │  │
│  │ ├─ auditMiddleware.Wrap() [Routes Audit log]     │  │
│  │ ├─ apiAuthMiddleware.Wrap() [Validates Bearer]   │  │
│  │ └─ roleMiddleware.RequireRole("admin")           │  │
│  └──────────────────────────────────────────────────┘  │
│              │                                         │
│              ▼                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │ HTTP Handler                                     │  │
│  │ (Parses request query/body, extracts context)    │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Application Layer (CQRS)                              │
│  ┌──────────────────────────────────────────────────┐  │
│  │ CQRS Command & Query Handlers                    │  │
│  │ (Executes logic, enforces domain validations)      │  │
│  │                                                  │  │
│  │ ├─ Uses Domain Models                            │  │
│  │ └─ Calls Repository Interfaces                   │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Domain Layer                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Domain Models & Core Logic                       │  │
│  │ ├─ User (Failed Attempts & Locks)                │  │
│  │ ├─ Client                                        │  │
│  │ ├─ AuthorizationCode                             │  │
│  │ ├─ AccessToken & RefreshToken                    │  │
│  │ ├─ Consent                                       │  │
│  │ ├─ Audit                                         │  │
│  │ └─ PKCE (S256 Validator)                         │  │
│  └──────────────────────────────────────────────────┘  │
│              │                                         │
│              ▼                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Repository Interfaces                            │  │
│  │ (Data access abstractions)                       │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Infrastructure Layer                                  │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Persistence & Persistence Wrapper                │  │
│  │ ├─ AsyncAuditRepository (Background channel)     │  │
│  │ └─ PostgreSQL Repositories (Prepared SQL)        │  │
│  └──────────────────────────────────────────────────┘  │
│              │                                         │
│              ▼                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Security Utilities                               │  │
│  │ ├─ PasswordHasher (bcrypt cost 14)               │  │
│  │ ├─ JWTSigner (asymmetric RS256 token helper)     │  │
│  │ ├─ JWKSProvider (manages/exposes RSA public key) │  │
│  │ └─ refresh_hash (crypto token generator)         │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
             │
             ▼
         PostgreSQL
          Database
```

---

## Dependency Injection Flow

The application boots in `cmd/api/main.go` using standard Go components:

```
main.go
  │
  ├─ Open Database Connection (sql.Open)
  │   └─ Optimize DB Connection Pool (SetMaxOpenConns: 100, SetMaxIdleConns: 50)
  │
  ├─ Initialize Repositories (PostgreSQL)
  │   ├─ UserRepository
  │   ├─ ClientRepository
  │   ├─ AuthorizationCodeRepository
  │   ├─ TokenRepository
  │   ├─ RefreshTokenRepository
  │   ├─ ConsentRepository
  │   └─ AuditRepository (Base postgres)
  │
  ├─ Initialize Async Audit Wrapper
  │   └─ AsyncAuditRepository (with buffer size 1000)
  │       └─ StartWorker() -> launches concurrent background writing thread
  │
  ├─ Initialize Infrastructure Security Services
  │   ├─ JWKSProvider (generates/rotates keys)
  │   ├─ SessionManager (HMAC-signed cookies)
  │   └─ JWTSigner (using JWKS keys)
  │
  ├─ Initialize CQRS Command Handlers
  │   ├─ AuthorizeHandler
  │   ├─ TokenHandler
  │   ├─ RefreshHandler
  │   ├─ LoginHandler
  │   ├─ ConsentHandler (command)
  │   ├─ RevokeHandler
  │   └─ RegisterClientHandler
  │
  ├─ Initialize CQRS Query Handlers
  │   ├─ AuditQueryHandler
  │   ├─ UserinfoHandler
  │   └─ ClientQueryHandler
  │
  ├─ Initialize HTTP Handlers
  │   ├─ AuthorizeHTTPHandler (handles authorize UI redirects)
  │   ├─ LoginHTTPHandler
  │   ├─ ConsentHTTPHandler
  │   ├─ TokenHTTPHandler (handles standard token operations)
  │   ├─ JWKSHTTPHandler
  │   ├─ AuditHTTPHandler
  │   ├─ OIDCHTTPHandler
  │   ├─ RevokeHTTPHandler
  │   ├─ RegistrationHTTPHandler
  │   ├─ UserinfoHTTPHandler
  │   └─ AdminHTTPHandler
  │
  ├─ Initialize HTTP Middlewares
  │   ├─ AuditMiddleware
  │   ├─ AuthMiddleware (Bearer extraction)
  │   ├─ RateLimiter (10,000 requests/minute)
  │   ├─ CORSMiddleware
  │   ├─ RequestLogger
  │   └─ RoleMiddleware (Admin validation)
  │
  ├─ Register Routes
  │   └─ Wrap endpoints in respective middlewares
  │
  └─ Start HTTP Server (http.ListenAndServe)
```

---

## Data Model Relationships

```
┌──────────────┐
│    User      │
│  (Lockout)   │
├──────────────┤
│ id (PK)      │
│ username     │
│ password     │
│ role         │
└──────┬───────┘
       │ (1)
       │
       ├──┬─────────────────────────────────────┐
       │  │                                     │
       │ (M)                                   (M)
      (M) │                                     │
       │  └──────────────┐                      │
       │                 │                      │
       ▼                 ▼                      ▼
┌─────────────────┐  ┌──────────┐         ┌─────────┐
│Authorization    │  │ Consent  │         │ Tokens  │
│Code             │  │          │         │         │
└─────────────────┘  └──────────┘         │ - AT    │
                                          │ - RT    │
                                          └─────────┘

┌──────────────┐
│   Client     │
├──────────────┤
│ id (PK)      │
│ client_id    │
│ secret       │
│ redirect_uri │
└──────┬───────┘
       │ (1)
       │
       ├──┬─────────────────────────┐
       │  │                         │
     (M) (M)                       (M)
       │  │                         │
       ▼  ▼                         ▼
   Auth  Consent              Tokens
   Code               
       
All entities record security audit records in the audit_logs table.
```

---

## CQRS Pattern Implementation

Separates reading operations (Queries) from state-changing actions (Commands) to allow scaling components independently:

```
User Write Request
     │
     ▼
┌─────────────────┐          ┌──────────────────────────────────┐
│ WRITE COMMAND   │          │ COMMAND HANDLER                  │
│ (e.g. Login)    │ ─────────► (Validates credentials, checks   │
└─────────────────┘          │  lockout times, writes audit log)│
                             └────────────────┬─────────────────┘
                                              │
                                              ▼
                                       Postgres database

─────────────────────────────────────────────────────────────────

User Read Request
     │
     ▼
┌─────────────────┐          ┌──────────────────────────────────┐
│ READ QUERY      │          │ QUERY HANDLER                    │
│ (e.g. Userinfo) │ ─────────► (Resolves profile payload by     │
└─────────────────┘          │  fetching from UserRepository)   │
                             └────────────────┬─────────────────┘
                                              │
                                              ▼
                                      Returned payload
```

---

## File Organization by Feature

### Authentication & Lockout Feature
```
domain/
  ├─ user.go (User lockout conditions)
  └─ audit.go (Logging audit records)

application/command/
  └─ login_handler.go (Lockout checking logic)

infrastructure/
  ├─ persistence/postgres/user_repo.go
  └─ security/password.go (bcrypt)

interfaces/http/
  ├─ handler_login.go
  └─ session.go (Cookie session storage)
```

### Authorization (PKCE & OIDC)
```
domain/
  ├─ auth_code.go
  ├─ client.go
  ├─ consent.go
  └─ pkce.go

application/command/
  ├─ authorize_handler.go
  └─ consent_handler.go

application/query/
  └─ userinfo_handler.go

interfaces/http/
  ├─ handler_authorize.go
  ├─ handler_consent.go
  ├─ handler_oidc.go
  └─ handler_userinfo.go
```

### Token Operations (Rotation & Revocation)
```
domain/
  └─ token.go

application/command/
  ├─ token_handler.go
  ├─ refresh_handler.go (Token rotation logic)
  └─ revoke_handler.go (Token revocation logic)

interfaces/http/
  ├─ handler_token.go
  └─ handler_revoke.go
```

---

## Technology Stack

```
Language & Version:
  └─ Go 1.26.1
  
API Library:
  └─ net/http (standard library)

Database Persistence:
  └─ PostgreSQL 15+
  └─ github.com/lib/pq (PostgreSQL driver)

Cryptography & Security:
  └─ golang.org/x/crypto/bcrypt (Password hashes)
  └─ crypto/rsa, crypto/sha256 (JWT RS256 Custom signature)

Deployment & Environment:
  └─ Docker
  └─ Docker Compose
  └─ Make (Task automation runner)
```
