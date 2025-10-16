package infrastructure

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// UserRepository is a PostgreSQL implementation of the domain.UserRepository.
type UserRepository struct {
	DB *infrastructure.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *infrastructure.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, is_email_verified, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.DB.Pool.QueryRow(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.IsEmailVerified,
		user.AvatarURL,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// GetUserByID retrieves a user from the database by their ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at, is_email_verified, deleted_at, avatar_url, deletion_due_at, last_login_at, is_deleted
		FROM users
		WHERE id = $1`

	err := r.DB.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsEmailVerified,
		&user.DeletedAt,
		&user.AvatarURL,
		&user.DeletionDueAt,
		&user.LastLoginAt,
		&user.IsDeleted,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user from the database by their email.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at, is_email_verified, deleted_at, avatar_url, deletion_due_at, last_login_at, is_deleted
		FROM users
		WHERE email = $1`

	err := r.DB.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsEmailVerified,
		&user.DeletedAt,
		&user.AvatarURL,
		&user.DeletionDueAt,
		&user.LastLoginAt,
		&user.IsDeleted,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user in the database.
func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, password_hash = $4, is_email_verified = $5, avatar_url = $6, deleted_at = $7, deletion_due_at = $8, last_login_at = $9, is_deleted = $10
		WHERE id = $1
		RETURNING updated_at`

	return r.DB.Pool.QueryRow(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.IsEmailVerified,
		user.AvatarURL,
		user.DeletedAt,
		user.DeletionDueAt,
		user.LastLoginAt,
		user.IsDeleted,
	).Scan(&user.UpdatedAt)
}

// DeleteUser deletes a user from the database.
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_deleted = true, deleted_at = now() WHERE id = $1`
	_, err := r.DB.Pool.Exec(ctx, query, id)
	return err
}
