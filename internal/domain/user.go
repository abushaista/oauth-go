package domain

import "context"

type User struct {
	ID       string
	Username string
	Password string
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
