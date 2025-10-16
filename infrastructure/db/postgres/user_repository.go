package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
	"github.com/jefersonprimer/chatear-backend/domain/repositories"
)

// PgxPoolIface defines the methods of pgxpool.Pool that userRepository uses
type PgxPoolIface interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// userRepository implements the UserRepository interface
type userRepository struct {
	pool PgxPoolIface
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool PgxPoolIface) repositories.UserRepository {
	return &userRepository{
		pool: pool,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, is_email_verified, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.IsEmailVerified,
		user.AvatarURL,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by ID from the database
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at, is_email_verified, 
		       deleted_at, avatar_url, deletion_due_at, last_login_at, is_deleted
		FROM users
		WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
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

// GetByEmail retrieves a user by email from the database
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at, is_email_verified, 
		       deleted_at, avatar_url, deletion_due_at, last_login_at, is_deleted
		FROM users
		WHERE email = $1`

	err := r.pool.QueryRow(ctx, query, email).Scan(
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

// Update updates a user in the database
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, password_hash = $4, is_email_verified = $5, 
		    avatar_url = $6, deleted_at = $7, deletion_due_at = $8, last_login_at = $9, is_deleted = $10
		WHERE id = $1
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
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

// Delete deletes a user from the database (soft delete)
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_deleted = true, deleted_at = now() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// GetDeletedUsers retrieves deleted users
func (r *userRepository) GetDeletedUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at, is_email_verified, 
		       deleted_at, avatar_url, deletion_due_at, last_login_at, is_deleted
		FROM users
		WHERE is_deleted = true
		ORDER BY deleted_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
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
		users = append(users, user)
	}

	return users, nil
}

// GetByEmailVerified retrieves users by email verification status
func (r *userRepository) GetByEmailVerified(ctx context.Context, verified bool, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at, is_email_verified, 
		       deleted_at, avatar_url, deletion_due_at, last_login_at, is_deleted
		FROM users
		WHERE is_email_verified = $1 AND is_deleted = false
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, verified, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
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
		users = append(users, user)
	}

	return users, nil
}

// SearchByName searches users by name
func (r *userRepository) SearchByName(ctx context.Context, name string, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at, is_email_verified, 
		       deleted_at, avatar_url, deletion_due_at, last_login_at, is_deleted
		FROM users
		WHERE name ILIKE $1 AND is_deleted = false
		ORDER BY name
		LIMIT $2 OFFSET $3`

	searchPattern := fmt.Sprintf("%%%s%%", name)
	rows, err := r.pool.Query(ctx, query, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
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
		users = append(users, user)
	}

	return users, nil
}
