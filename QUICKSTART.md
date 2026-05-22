# Quick Start Guide

This guide gets you up and running with the OAuth 2.1 & OIDC Server.

## 🚀 Setup Development Environment

### Option 1: Using Docker (Recommended)

1. **Start the environment (PostgreSQL + Server + Adminer)**:
   ```bash
   make docker-up
   ```
   *Note: Database migrations and test seeds are automatically executed when the Postgres container starts.*

2. **Accessing the services**:
   - **OAuth Server**: `http://localhost:8080`
   - **Adminer (Database UI)**: `http://localhost:8081`

3. **Stop the environment**:
   ```bash
   make docker-down
   ```

### Option 2: Local Server + Dockerized DB (Best for code edits)

1. **Start the database container**:
   ```bash
   make docker-db
   ```

2. **Start the Go server locally**:
   ```bash
   make dev
   ```
   *This loads environment variables from Makefile and starts the hot-reloading server connected to Postgres on port 8080.*

---

## 🎯 API Testing

Here are the curl commands to verify and interact with the endpoints.

### 1. Health Check
```bash
curl http://localhost:8080/health
# Response: OK
```

### 2. OpenID Connect Discovery
Fetch the provider metadata:
```bash
curl http://localhost:8080/.well-known/openid-configuration
```

### 3. Retrieve JWKS
Get the public keys used for signature verification:
```bash
curl http://localhost:8080/.well-known/jwks.json
```

### 4. Direct User Login
Authenticate credentials to start a session:
```bash
curl -i -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```
*Note: Successful login sets an HMAC-signed cookie `session_id`.*

### 5. Dynamic Client Registration (RFC 7591)
Register a new client application dynamically:
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "My Client App",
    "redirect_uris": ["http://localhost:3000/callback"]
  }'
```
*Response includes client ID and client secret.*

### 6. Token Exchange (Client Credentials Flow)
Get a token directly using the client's credentials:
```bash
curl -X POST http://localhost:8080/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=demo-app&client_secret=demo-secret"
```

### 7. Token Revocation (RFC 7009)
Revoke an issued access or refresh token:
```bash
curl -X POST http://localhost:8080/oauth/revoke \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "token=YOUR_ACCESS_OR_REFRESH_TOKEN&client_id=demo-app&client_secret=demo-secret"
```

### 8. Userinfo Profile Endpoint
Fetch user metadata using a bearer token:
```bash
curl http://localhost:8080/userinfo \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 9. Audits API (Admin Only)
Requires admin role bearer token. Run:
```bash
curl "http://localhost:8080/audits?user_id=test-user-001" \
  -H "Authorization: Bearer YOUR_ADMIN_ACCESS_TOKEN"
```

---

## 🗄️ Database Management

### Using Adminer UI
- **URL**: `http://localhost:8081`
- **System**: PostgreSQL
- **Server**: `postgres` (or `localhost` if connecting from outside Docker)
- **Username**: `oauth_user`
- **Password**: `oauth_password`
- **Database**: `oauth_db`

### Using psql CLI (Local)
Connect directly via shell:
```bash
# Connect to container database
psql -h localhost -U oauth_user -d oauth_db

# List tables
\dt

# View seed user with role admin and lockout stats
SELECT id, username, role, failed_login_attempts FROM users;

# Exit
\q
```

---

## 🚀 Next Steps

- [ ] Write integration tests for the dynamic client registration flow.
- [ ] Connect a client application (e.g. NextJS or Vue app) to the server.
- [ ] Implement database-backed client creation screens.
