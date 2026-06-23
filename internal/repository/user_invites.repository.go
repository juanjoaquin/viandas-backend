package repository

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveUserInvite(ctx context.Context, email, role, token, invitedBy string, expiresAt time.Time) (*entity.UserInvite, error) {
	var invite entity.UserInvite
	err := r.db.GetContext(ctx, &invite, `
		INSERT INTO user_invites (email, role, token, invited_by, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`, email, role, token, invitedBy, expiresAt)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *repo) GetUserInviteByToken(ctx context.Context, token string) (*entity.UserInvite, error) {
	var invite entity.UserInvite
	err := r.db.GetContext(ctx, &invite, `SELECT * FROM user_invites WHERE token = $1`, token)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *repo) AcceptUserInvite(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_invites
		SET accepted_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
