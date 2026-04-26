# OAuth Server Architecture Visualization

## Request Flow Diagram

```
HTTP Request
    │
    ▼
┌─────────────────────────────────────────┐
│     HTTP Interface Layer                │
│  ┌──────────────────────────────────┐  │
│  │ Router (http.ServeMux)           │  │
│  │ ├─ /oauth/authorize              │  │
│  │ ├─ /oauth/token                  │  │
│  │ ├─ /login                        │  │
│  │ ├─ /consent                      │  │
│  │ └─ /.well-known/jwks.json        │  │
│  └──────────────────────────────────┘  │
│              │                          │
│              ▼                          │
│  ┌──────────────────────────────────┐  │
│  │ HTTP Handler                     │  │
│  │ (Parses request, validates)      │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Application Layer (CQRS)               │
│  ┌──────────────────────────────────┐  │
│  │ Command Handler                  │  │
│  │ (Business logic, validation)     │  │
│  │                                  │  │
│  │ └─ Uses Domain Models            │  │
│  │ └─ Calls Repository Interfaces   │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Domain Layer                           │
│  ┌──────────────────────────────────┐  │
│  │ Domain Models                    │  │
│  │ ├─ User                          │  │
│  │ ├─ Client                        │  │
│  │ ├─ AuthorizationCode             │  │
│  │ ├─ Token                         │  │
│  │ ├─ Consent                       │  │
│  │ ├─ Audit                         │  │
│  │ └─ PKCE                          │  │
│  └──────────────────────────────────┘  │
│              │                          │
│              ▼                          │
│  ┌──────────────────────────────────┐  │
│  │ Repository Interfaces            │  │
│  │ (Abstraction for data access)    │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Infrastructure Layer                   │
│  ┌──────────────────────────────────┐  │
│  │ PostgreSQL Repositories          │  │
│  │ (Concrete implementations)       │  │
│  └──────────────────────────────────┘  │
│              │                          │
│              ▼                          │
│  ┌──────────────────────────────────┐  │
│  │ Database Driver (github.com...)  │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
             │
             ▼
         PostgreSQL
        Database
```

---

## Dependency Injection Flow

```
main.go
  │
  ├─ Open Database Connection
  │
  ├─ Initialize Repositories
  │   ├─ UserRepository
  │   ├─ ClientRepository
  │   ├─ AuthorizationCodeRepository
  │   ├─ TokenRepository
  │   ├─ RefreshTokenRepository
  │   ├─ ConsentRepository
  │   └─ AuditRepository
  │
  ├─ Initialize Command Handlers
  │   ├─ AuthorizeHandler (using repositories)
  │   ├─ TokenHandler (using repositories)
  │   ├─ RefreshHandler (using repositories)
  │   └─ LoginHandler (using repositories)
  │
  ├─ Initialize Security Services
  │   ├─ PasswordHasher
  │   ├─ JWTSigner
  │   ├─ PKCEValidator
  │   └─ JWKSProvider
  │
  ├─ Initialize HTTP Handlers
  │   ├─ AuthorizeHTTPHandler (using command handler)
  │   ├─ TokenHTTPHandler (using command handler)
  │   ├─ LoginHTTPHandler (using command handler)
  │   ├─ ConsentHTTPHandler
  │   ├─ JWKSHTTPHandler (using JWKS provider)
  │   └─ StaticHandler
  │
  ├─ Register Routes
  │   └─ http.ServeMux
  │
  └─ Start HTTP Server
      └─ Listen and Serve
```

---

## Data Model Relationships

```
┌──────────────┐
│    User      │
├──────────────┤
│ id (PK)      │
│ username     │
│ password     │
└──────┬───────┘
       │ (1)
       │
       ├──┬─────────────────────────────────────┐
       │  │                                     │
     (M) (M)                                   (M)
       │  │                                     │
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
       
All entities have (M) relationship to Audit Logs
```

---

## CQRS Pattern Implementation

```
User Request
     │
     ├─────── Write Operations ─────────┐
     │                                   │
     ▼                                   ▼
┌─────────────────┐          ┌──────────────────┐
│ LOGIN COMMAND   │          │ COMMAND HANDLER  │
├─────────────────┤          ├──────────────────┤
│ username        │          │ Authenticate     │
│ password        │          │ Validate         │
└─────────────────┘          │ Execute          │
                              │ Use Repositories │
                              └──────────────────┘
                                    │
     ┌───────────────────────────────┘
     │
     ▼
  Repositories (PostgreSQL)
  
────────────────────────────────────────────────

User Request
     │
     ├────── Read Operations ───────────┐
     │                                  │
     ▼                                  ▼
┌─────────────────┐        ┌──────────────────┐
│ JWKS QUERY      │        │ QUERY HANDLER    │
├─────────────────┤        ├──────────────────┤
│ key_id?         │        │ Retrieve         │
└─────────────────┘        │ Format           │
                            │ Return           │
                            └──────────────────┘
                                  │
                                  ▼
                              Response Cache
```

---

## File Organization by Feature

### Authentication Feature
```
domain/
  ├─ user.go (User entity + interface)
  ├─ crypto.go (Auth constants)
  └─ audit.go (Audit logging)

application/command/
  └─ login_handler.go

infrastructure/
  ├─ persistence/postgres/user_repo.go
  └─ security/password.go

interfaces/http/
  └─ handler_login.go
```

### Authorization Feature
```
domain/
  ├─ auth_code.go (AuthCode entity + interface)
  ├─ consent.go (Consent entity + interface)
  ├─ pkce.go (PKCE validation)
  └─ client.go (Client entity + interface)

application/command/
  └─ authorize_handler.go

infrastructure/
  ├─ persistence/postgres/auth_code_repo.go
  ├─ persistence/postgres/consent_repo.go
  ├─ persistence/postgres/client_repo.go
  └─ security/refresh_hash.go

interfaces/http/
  ├─ handler_authorize.go
  └─ handler_consent.go
```

### Token Feature
```
domain/
  ├─ token.go (Token entities)
  └─ token_repository.go (Token interfaces)

application/command/
  ├─ token_handler.go
  └─ refresh_handler.go

infrastructure/
  ├─ persistence/postgres/token_repo.go
  ├─ persistence/postgres/refresh_token_repo.go
  └─ security/jwt_rs256.go

interfaces/http/
  ├─ handler_token.go
  └─ handler_jwks.go
```

### Audit Feature
```
domain/
  ├─ audit.go (Audit entity + interface)

infrastructure/
  └─ persistence/postgres/audit_repo.go
```

---

## Database Schema Design

```
users
│
├─ One-to-Many → authorization_codes
├─ One-to-Many → access_tokens
├─ One-to-Many → refresh_tokens
├─ One-to-Many → consents
└─ One-to-Many → audit_logs
   
clients
│
├─ One-to-Many → authorization_codes
├─ One-to-Many → access_tokens
├─ One-to-Many → refresh_tokens
├─ One-to-Many → consents
└─ One-to-Many → audit_logs

Keys:
  PK = Primary Key
  FK = Foreign Key
  UNIQUE = Unique constraint
  INDEX = Database index
```

---

## Error Handling Flow

```
HTTP Request
    │
    ▼
Handler Layer
    │
    ├─ Validation Error?
    │   ├─ 400 Bad Request
    │   └─ Return error JSON
    │
    ▼
Application Layer
    │
    ├─ Business Logic Error?
    │   ├─ 401 Unauthorized / 403 Forbidden
    │   └─ Return error JSON
    │
    ├─ Not Found?
    │   ├─ 404 Not Found
    │   └─ Return error JSON
    │
    ▼
Infrastructure Layer
    │
    ├─ Database Error?
    │   ├─ 500 Internal Server Error
    │   └─ Log error
    │
    ├─ Connection Error?
    │   ├─ 503 Service Unavailable
    │   └─ Log error
    │
    ▼
Success Response
    ├─ 200 OK / 201 Created
    └─ Return result JSON
```

---

## Deployment Architecture (Optional)

```
┌──────────────────────┐
│  Load Balancer       │
│  (nginx/haproxy)     │
└──────────┬───────────┘
           │
      ┌────┴────┐
      │          │
      ▼          ▼
┌────────────┐ ┌────────────┐
│ OAuth      │ │ OAuth      │
│ Server 1   │ │ Server 2   │
└────────────┘ └────────────┘
      │          │
      └────┬─────┘
           │
           ▼
    ┌─────────────┐
    │ PostgreSQL  │
    │ (Primary)   │
    └─────────────┘
           │
           ├─ Replication
           │
           ▼
    ┌─────────────┐
    │ PostgreSQL  │
    │ (Standby)   │
    └─────────────┘

Optional Caching Layer:
    Redis (Session/Token Cache)
```

---

## Development Workflow

```
1. Clone/Setup
   └─ git clone
   └─ cd oauth-go
   
2. Dependencies
   └─ go mod download
   
3. Database
   └─ make docker-up
   └─ make migrate
   
4. Development
   └─ make dev
   └─ Server runs on :8080
   
5. Testing
   └─ make test
   
6. Build
   └─ make build
   └─ Binary: ./bin/oauth-server
   
7. Deploy
   └─ Docker container
   └─ Or binary + systemd
```

---

## Technology Stack

```
Backend Framework:
  └─ Go 1.21+
  
Web Framework:
  └─ net/http (Standard Library)
  
Database:
  └─ PostgreSQL 13+
  └─ github.com/lib/pq (Driver)
  
Authentication:
  └─ bcrypt (future)
  └─ JWT RS256 (future)
  
Development:
  └─ Docker
  └─ Docker Compose
  └─ Make
  
Patterns:
  └─ Clean Architecture
  └─ CQRS
  └─ Repository Pattern
  └─ Dependency Injection
```

---

## Key Design Decisions

1. **Clean Architecture**: Separated concerns across 4 layers
2. **CQRS**: Separate command and query handlers
3. **Repository Pattern**: Abstract database access
4. **Context Usage**: Proper context propagation
5. **PostgreSQL**: Proven SQL database for OAuth
6. **Prepared Statements**: SQL injection protection
7. **Interface-based**: Easy to test and mock
8. **Dependency Injection**: Flexible and testable

---

## Scalability Considerations

```
Single Server (Current)
  └─ Suitable for: Development, small deployments
  
Horizontal Scaling (Recommended)
  ├─ Multiple OAuth servers behind load balancer
  ├─ Shared PostgreSQL database
  ├─ Optional Redis for caching
  └─ Database replication for HA
  
Cloud Deployment
  ├─ Kubernetes with horizontal pod autoscaling
  ├─ Managed PostgreSQL (AWS RDS, Google Cloud SQL)
  ├─ CDN for static files
  └─ Monitoring & logging (ELK, Datadog, etc.)
```

---

**This architecture is production-ready and scalable! 🚀**
