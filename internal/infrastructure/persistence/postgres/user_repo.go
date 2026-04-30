package postgres

import (
	"context"
	"database/sql"

	"github.com/abushaista/oauth-go/internal/domain"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, username, password, role, failed_login_attempts, locked_until FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.FailedLoginAttempts, &user.LockedUntil)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, username, password, role, failed_login_attempts, locked_until FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.FailedLoginAttempts, &user.LockedUntil)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Save creates or updates a user
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, username, password, role, failed_login_attempts, locked_until) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 ON CONFLICT (id) DO UPDATE SET 
			username = $2, 
			password = $3, 
			role = $4,
			failed_login_attempts = $5, 
			locked_until = $6`,
		user.ID, user.Username, user.Password, user.Role, user.FailedLoginAttempts, user.LockedUntil,
	)
	return err
}
