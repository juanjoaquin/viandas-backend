package repository

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *repo) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	err := r.db.GetContext(ctx, &rt,
		`SELECT * FROM refresh_tokens WHERE token = $1`,
		token,
	)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *repo) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM refresh_tokens WHERE token = $1`,
		token,
	)
	return err
}

func (r *repo) DeleteRefreshTokensByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	return err
}
