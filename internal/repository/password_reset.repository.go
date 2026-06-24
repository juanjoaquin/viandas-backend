package repository

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SavePasswordReset(ctx context.Context, userID, token string, expiresAt time.Time) (*entity.PasswordReset, error) {
	var pr entity.PasswordReset
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO password_resets (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING *`,
		userID, token, expiresAt,
	).StructScan(&pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *repo) GetPasswordResetByToken(ctx context.Context, token string) (*entity.PasswordReset, error) {
	var pr entity.PasswordReset
	err := r.db.GetContext(ctx, &pr,
		`SELECT * FROM password_resets WHERE token = $1`,
		token,
	)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *repo) MarkPasswordResetUsed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE password_resets SET used_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}

func (r *repo) UpdateUserPassword(ctx context.Context, userID, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		passwordHash, userID,
	)
	return err
}
