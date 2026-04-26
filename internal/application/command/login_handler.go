package command

import (
	"context"
	"fmt"
	"time"

	"github.com/abushaista/oauth-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler authenticates users and manages login state
type LoginHandler struct {
	userRepo  domain.UserRepository
	auditRepo domain.AuditRepository
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(userRepo domain.UserRepository, auditRepo domain.AuditRepository) *LoginHandler {
	return &LoginHandler{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

// Handle processes a login command
func (h *LoginHandler) Handle(ctx context.Context, cmd Command) (interface{}, error) {
	loginCmd, ok := cmd.(*LoginCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	// Validate input
	if loginCmd.Username == "" || loginCmd.Password == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	// Find user by username
	user, err := h.userRepo.FindByUsername(ctx, loginCmd.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCmd.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Make Audit log entry
	auditEntry := &domain.Audit{
		ID:        generateRandomToken(),
		UserID:    user.ID,
		ClientID:  "system", // Client is not directly applicable during login for now unless context supplies it
		Action:    "LOGIN_SUCCESS",
		Details:   "User successfully authenticated.",
		IPAddress: "", // Could be passed in from context
		CreatedAt: time.Now(),
	}
	// Best effort, don't block login if audit fails
	_ = h.auditRepo.Create(ctx, auditEntry)

	return map[string]interface{}{
		"user_id": user.ID,
		"success": true,
	}, nil
}
