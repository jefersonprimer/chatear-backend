package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// PostgresRefreshTokenRepository implements domain.RefreshTokenRepository for PostgreSQL.
type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

// NewPostgresRefreshTokenRepository creates a new PostgresRefreshTokenRepository.
func NewPostgresRefreshTokenRepository(db *sql.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

// CreateRefreshToken creates a new refresh token in the database.
func (r *PostgresRefreshTokenRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx,
		query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

// GetRefreshTokenByToken finds a refresh token by its token string.
func (r *PostgresRefreshTokenRepository) GetRefreshTokenByToken(ctx context.Context, tokenString string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	var revokedAt sql.NullTime

	query := `SELECT id, user_id, token, expires_at, created_at, revoked_at FROM refresh_tokens WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, tokenString).Scan(
		&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt, &revokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}

	return &token, nil
}

// UpdateRefreshToken updates a refresh token in the database.
func (r *PostgresRefreshTokenRepository) UpdateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, token.RevokedAt, token.ID)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}
	return nil
}

// RevokeRefreshToken revokes a refresh token by its ID.
func (r *PostgresRefreshTokenRepository) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a given user ID.
func (r *PostgresRefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user refresh tokens: %w", err)
	}
	return nil
}
