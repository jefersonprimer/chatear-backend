package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
	"github.com/jefersonprimer/chatear-backend/domain/repositories"
)

// magicLinkRepository implements the MagicLinkRepository interface
type magicLinkRepository struct {
	adapter *Adapter
}

// NewMagicLinkRepository creates a new magic link repository
func NewMagicLinkRepository(adapter *Adapter) repositories.MagicLinkRepository {
	return &magicLinkRepository{
		adapter: adapter,
	}
}

// Create creates a new magic link
func (r *magicLinkRepository) Create(ctx context.Context, magicLink *entities.MagicLink) error {
	query := `
		INSERT INTO magic_links (id, user_id, token, expires_at, used, created_at, type, used_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.adapter.Pool.Exec(ctx, query,
		magicLink.ID,
		magicLink.UserID,
		magicLink.Token,
		magicLink.ExpiresAt,
		magicLink.Used,
		magicLink.CreatedAt,
		magicLink.Type,
		magicLink.UsedAt,
		magicLink.IsActive,
	)

	return err
}

// GetByToken retrieves a magic link by token
func (r *magicLinkRepository) GetByToken(ctx context.Context, token string) (*entities.MagicLink, error) {
	magicLink := &entities.MagicLink{}
	query := `
		SELECT id, user_id, token, expires_at, used, created_at, type, used_at, is_active
		FROM magic_links
		WHERE token = $1`

	err := r.adapter.Pool.QueryRow(ctx, query, token).Scan(
		&magicLink.ID,
		&magicLink.UserID,
		&magicLink.Token,
		&magicLink.ExpiresAt,
		&magicLink.Used,
		&magicLink.CreatedAt,
		&magicLink.Type,
		&magicLink.UsedAt,
		&magicLink.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return magicLink, nil
}

// GetByUserID retrieves magic links by user ID
func (r *magicLinkRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.MagicLink, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at, type, used_at, is_active
		FROM magic_links
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.adapter.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMagicLinks(rows)
}

// GetByUserIDAndType retrieves magic links by user ID and type
func (r *magicLinkRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, linkType entities.MagicLinkType, limit, offset int) ([]*entities.MagicLink, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at, type, used_at, is_active
		FROM magic_links
		WHERE user_id = $1 AND type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.adapter.Pool.Query(ctx, query, userID, linkType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMagicLinks(rows)
}

// GetActiveByUserIDAndType retrieves active magic links by user ID and type
func (r *magicLinkRepository) GetActiveByUserIDAndType(ctx context.Context, userID uuid.UUID, linkType entities.MagicLinkType) ([]*entities.MagicLink, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at, type, used_at, is_active
		FROM magic_links
		WHERE user_id = $1 AND type = $2 AND is_active = true AND used = false AND expires_at > now()
		ORDER BY created_at DESC`

	rows, err := r.adapter.Pool.Query(ctx, query, userID, linkType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMagicLinks(rows)
}

// Update updates a magic link
func (r *magicLinkRepository) Update(ctx context.Context, magicLink *entities.MagicLink) error {
	query := `
		UPDATE magic_links
		SET user_id = $2, token = $3, expires_at = $4, used = $5, type = $6, used_at = $7, is_active = $8
		WHERE id = $1`

	_, err := r.adapter.Pool.Exec(ctx, query,
		magicLink.ID,
		magicLink.UserID,
		magicLink.Token,
		magicLink.ExpiresAt,
		magicLink.Used,
		magicLink.Type,
		magicLink.UsedAt,
		magicLink.IsActive,
	)

	return err
}

// Delete deletes a magic link by ID
func (r *magicLinkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM magic_links WHERE id = $1`
	_, err := r.adapter.Pool.Exec(ctx, query, id)
	return err
}

// DeleteExpired deletes expired magic links
func (r *magicLinkRepository) DeleteExpired(ctx context.Context, olderThan time.Time) error {
	query := `DELETE FROM magic_links WHERE expires_at < $1`
	_, err := r.adapter.Pool.Exec(ctx, query, olderThan)
	return err
}

// RevokeByUserIDAndType revokes magic links by user ID and type
func (r *magicLinkRepository) RevokeByUserIDAndType(ctx context.Context, userID uuid.UUID, linkType entities.MagicLinkType) error {
	query := `UPDATE magic_links SET is_active = false WHERE user_id = $1 AND type = $2`
	_, err := r.adapter.Pool.Exec(ctx, query, userID, linkType)
	return err
}

// scanMagicLinks scans rows into magic link entities
func (r *magicLinkRepository) scanMagicLinks(rows interface{}) ([]*entities.MagicLink, error) {
	// This would need to be implemented based on the actual pgx rows interface
	// For now, returning empty slice as placeholder
	return []*entities.MagicLink{}, nil
}
