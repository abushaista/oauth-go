package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/abushaista/oauth-go/internal/application/command"
	"github.com/abushaista/oauth-go/internal/application/query"
	"github.com/abushaista/oauth-go/internal/infrastructure/persistence"
	"github.com/abushaista/oauth-go/internal/infrastructure/persistence/postgres"
	"github.com/abushaista/oauth-go/internal/infrastructure/security"
	handlers "github.com/abushaista/oauth-go/internal/interfaces/http"
)

func main() {
	// Load configuration from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/oauth_db?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Optimize DB connection pool
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to database (Pool: 100/50)")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	clientRepo := postgres.NewClientRepository(db)
	authCodeRepo := postgres.NewAuthorizationCodeRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	consentRepo := postgres.NewConsentRepository(db)
	
	// Use AsyncAuditRepository to offload DB writes
	baseAuditRepo := postgres.NewAuditRepository(db)
	auditRepo := persistence.NewAsyncAuditRepository(baseAuditRepo, 1000)
	auditRepo.StartWorker(context.Background())

	// Initialize infrastructure services
	jwksProvider := security.NewJWKSProvider()
	sessionManager := handlers.NewSessionManager()
	jwtSigner := security.NewJWTSigner(jwksProvider)

	// Initialize command handlers
	authorizeHandler := command.NewAuthorizeHandler(authCodeRepo, clientRepo, consentRepo, auditRepo)
	tokenHandler := command.NewTokenHandler(authCodeRepo, tokenRepo, refreshTokenRepo, clientRepo, auditRepo, userRepo, jwtSigner)
	refreshHandler := command.NewRefreshHandler(refreshTokenRepo, tokenRepo, clientRepo, auditRepo)
	loginHandler := command.NewLoginHandler(userRepo, auditRepo)
	consentCommandHandler := command.NewConsentHandler(consentRepo)
	revokeHandler := command.NewRevokeHandler(tokenRepo, refreshTokenRepo, clientRepo, auditRepo)
	registrationHandler := command.NewRegisterClientHandler(clientRepo, auditRepo)

	// Initialize query handlers
	auditQueryHandler := query.NewAuditQueryHandler(auditRepo)
	userinfoQueryHandler := query.NewUserinfoHandler(userRepo)
	clientQueryHandler := query.NewClientQueryHandler(clientRepo)



	// Initialize HTTP handlers
	authHTTPHandler := handlers.NewAuthorizeHandler(authorizeHandler, sessionManager)
	loginHTTPHandler := handlers.NewLoginHandler(loginHandler, sessionManager)
	consentHTTPHandler := handlers.NewConsentHandler(consentCommandHandler, sessionManager)
	tokenHTTPHandler := handlers.NewTokenHandler(tokenHandler, refreshHandler)
	jwksHTTPHandler := handlers.NewJWKSHandler(jwksProvider)
	auditHTTPHandler := handlers.NewAuditHandler(auditQueryHandler)
	oidcHTTPHandler := handlers.NewOIDCHandler()
	revokeHTTPHandler := handlers.NewRevokeHandler(revokeHandler)
	registrationHTTPHandler := handlers.NewRegistrationHandler(registrationHandler)
	userinfoHTTPHandler := handlers.NewUserinfoHandler(userinfoQueryHandler)
	adminHTTPHandler := handlers.NewAdminHandler(clientQueryHandler, auditQueryHandler)

	auditMiddleware := handlers.NewAuditMiddleware(auditRepo)
	apiAuthMiddleware := handlers.NewAuthMiddleware(tokenRepo)

	// Initialize global middlewares
	// Relaxed rate limit for high TPS (e.g., 10k requests per minute per IP)
	rateLimiter := handlers.NewRateLimiter(10000, 1*time.Minute) 
	corsMiddleware := handlers.NewCORSMiddleware()
	requestLogger := handlers.NewRequestLogger()

	roleMiddleware := handlers.NewRoleMiddleware(userRepo)

	// Register routes
	mux := http.NewServeMux()

	// OAuth endpoints (Wrapped with Middleware)
	mux.Handle("/oauth/authorize", auditMiddleware.Wrap(authHTTPHandler))
	mux.Handle("/oauth/token", auditMiddleware.Wrap(tokenHTTPHandler))
	mux.Handle("/oauth/revoke", auditMiddleware.Wrap(revokeHTTPHandler))
	mux.Handle("/.well-known/jwks.json", jwksHTTPHandler)
	mux.Handle("/.well-known/openid-configuration", oidcHTTPHandler)

	// User endpoints
	mux.Handle("/login", auditMiddleware.Wrap(loginHTTPHandler))
	mux.Handle("/consent", auditMiddleware.Wrap(consentHTTPHandler))
	
	// API Endpoints
	mux.Handle("/audits", apiAuthMiddleware.Wrap(roleMiddleware.RequireRole("admin")(auditHTTPHandler)))
	mux.Handle("/userinfo", apiAuthMiddleware.Wrap(userinfoHTTPHandler))
	mux.Handle("/register", registrationHTTPHandler)
	mux.Handle("/admin/api", apiAuthMiddleware.Wrap(roleMiddleware.RequireRole("admin")(adminHTTPHandler))) // Protected admin API

	// UI endpoints
	staticHandler := handlers.NewStaticHandler("web")
	mux.Handle("/ui/", http.StripPrefix("/ui/", staticHandler))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Start server with global middleware chain:
	// Request Logger → CORS → Rate Limiter → Router
	var handler http.Handler = mux
	handler = rateLimiter.Wrap(handler)
	handler = corsMiddleware.Wrap(handler)
	handler = requestLogger.Wrap(handler)

	log.Printf("Starting OAuth server on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
