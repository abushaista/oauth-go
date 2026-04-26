package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ErrInvalidSession = errors.New("invalid session")
	ErrExpiredSession = errors.New("session expired")
)

type SessionManager struct {
	secret []byte
}

func NewSessionManager() *SessionManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production"
	}
	return &SessionManager{
		secret: []byte(secret),
	}
}

// CreateSession generate a signed JWT for the user session
func (s *SessionManager) CreateSession(userID string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	payloadJSON, err := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadJSON)

	signingInput := header + "." + payload
	signature := s.sign(signingInput)

	return signingInput + "." + signature, nil
}

// VerifySession validates the session token and returns the user ID
func (s *SessionManager) VerifySession(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidSession
	}

	headerB64, payloadB64, signatureB64 := parts[0], parts[1], parts[2]

	// Verify signature
	signingInput := headerB64 + "." + payloadB64
	expectedSignature := s.sign(signingInput)

	if !hmac.Equal([]byte(signatureB64), []byte(expectedSignature)) {
		return "", ErrInvalidSession
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return "", ErrInvalidSession
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", ErrInvalidSession
	}

	// Check expiration
	exp, ok := payload["exp"].(float64)
	if !ok {
		return "", ErrInvalidSession
	}
	if time.Now().Unix() > int64(exp) {
		return "", ErrExpiredSession
	}

	userID, ok := payload["user_id"].(string)
	if !ok {
		return "", ErrInvalidSession
	}

	return userID, nil
}

func (s *SessionManager) sign(input string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// SetSessionCookie sets the cookie on the response
func (s *SessionManager) SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true for production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24,
	})
}
