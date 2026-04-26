package main

import (
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

	log.Println("Connected to database")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	clientRepo := postgres.NewClientRepository(db)
	authCodeRepo := postgres.NewAuthorizationCodeRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	consentRepo := postgres.NewConsentRepository(db)
	auditRepo := postgres.NewAuditRepository(db)

	// Initialize command handlers
	authorizeHandler := command.NewAuthorizeHandler(authCodeRepo, clientRepo, consentRepo, auditRepo)
	tokenHandler := command.NewTokenHandler(authCodeRepo, tokenRepo, refreshTokenRepo, clientRepo, auditRepo)
	refreshHandler := command.NewRefreshHandler(refreshTokenRepo, tokenRepo, clientRepo, auditRepo)
	loginHandler := command.NewLoginHandler(userRepo, auditRepo)
	consentCommandHandler := command.NewConsentHandler(consentRepo)

	// Initialize query handlers
	auditQueryHandler := query.NewAuditQueryHandler(auditRepo)

	// Initialize infrastructure services
	jwksProvider := security.NewJWKSProvider()
	sessionManager := handlers.NewSessionManager()

	// Initialize HTTP handlers
	authHTTPHandler := handlers.NewAuthorizeHandler(authorizeHandler, sessionManager)
	loginHTTPHandler := handlers.NewLoginHandler(loginHandler, sessionManager)
	consentHTTPHandler := handlers.NewConsentHandler(consentCommandHandler, sessionManager)
	tokenHTTPHandler := handlers.NewTokenHandler(tokenHandler, refreshHandler)
	jwksHTTPHandler := handlers.NewJWKSHandler(jwksProvider)
	auditHTTPHandler := handlers.NewAuditHandler(auditQueryHandler)
	auditMiddleware := handlers.NewAuditMiddleware(auditRepo)

	// Initialize global middlewares
	rateLimiter := handlers.NewRateLimiter(60, 1*time.Minute) // 60 req/min per IP
	corsMiddleware := handlers.NewCORSMiddleware()
	requestLogger := handlers.NewRequestLogger()

	// Register routes
	mux := http.NewServeMux()

	// OAuth endpoints (Wrapped with Middleware)
	mux.Handle("/oauth/authorize", auditMiddleware.Wrap(authHTTPHandler))
	mux.Handle("/oauth/token", auditMiddleware.Wrap(tokenHTTPHandler))
	mux.Handle("/.well-known/jwks.json", jwksHTTPHandler)

	// User endpoints
	mux.Handle("/login", auditMiddleware.Wrap(loginHTTPHandler))
	mux.Handle("/consent", auditMiddleware.Wrap(consentHTTPHandler))
	
	// API Endpoints
	mux.Handle("/audits", auditHTTPHandler)

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
