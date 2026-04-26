package http

import (
	"net/http"
	"os"
	"strings"
)

// CORSMiddleware handles Cross-Origin Resource Sharing headers
type CORSMiddleware struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
	maxAge         string
}

// NewCORSMiddleware creates a CORS middleware from environment or defaults
func NewCORSMiddleware() *CORSMiddleware {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000,http://localhost:3001,http://localhost:8080"
	}

	return &CORSMiddleware{
		allowedOrigins: strings.Split(origins, ","),
		allowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		allowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		maxAge:         "86400", // 24 hours
	}
}

func (c *CORSMiddleware) isOriginAllowed(origin string) bool {
	for _, allowed := range c.allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

// Wrap applies CORS headers to the response
func (c *CORSMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && c.isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.allowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.allowedHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", c.maxAge)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
