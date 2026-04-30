package domain

import (
	"context"
	"time"
)

type User struct {
	ID                  string
	Username            string
	Password            string
	Role                string
	FailedLoginAttempts int
	LockedUntil         *time.Time
}

// IsLocked checks if the user account is currently locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// ResetLoginAttempts resets the failed login counter and unlock time
func (u *User) ResetLoginAttempts() {
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// FindByUsername retrieves a user by username
	FindByUsername(ctx context.Context, username string) (*User, error)

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id string) (*User, error)

	// Save creates or updates a user
	Save(ctx context.Context, user *User) error
}
