package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/juanjoaquin/viandas-backend/encryption"
)

var (
	ErrPasswordResetTokenInvalid = errors.New("invalid or expired reset token")
	ErrPasswordResetTokenUsed    = errors.New("reset token already used")
)

func (s *serv) RequestPasswordReset(ctx context.Context, email string) error {
	email = normalizeEmail(email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// No revelar si el email existe o no
		return nil
	}

	token, err := generateSecureToken()
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	if _, err := s.repo.SavePasswordReset(ctx, user.ID, token, expiresAt); err != nil {
		return err
	}

	resetURL := s.passwordResetURL(token)
	return s.mailer.SendPasswordReset(ctx, email, resetURL)
}

func (s *serv) ResetPassword(ctx context.Context, token, newPassword string) error {
	pr, err := s.repo.GetPasswordResetByToken(ctx, strings.TrimSpace(token))
	if err != nil {
		return ErrPasswordResetTokenInvalid
	}

	if pr.UsedAt.Valid {
		return ErrPasswordResetTokenUsed
	}

	if time.Now().After(pr.ExpiresAt) {
		return ErrPasswordResetTokenInvalid
	}

	hashed, err := encryption.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateUserPassword(ctx, pr.UserID, hashed); err != nil {
		return err
	}

	return s.repo.MarkPasswordResetUsed(ctx, pr.ID)
}

func (s *serv) passwordResetURL(token string) string {
	baseURL := strings.TrimRight(s.settings.FrontendURL, "/")
	return fmt.Sprintf("%s/reset-password?token=%s", baseURL, url.QueryEscape(token))
}

func generateSecureToken() (string, error) {
	return generateInviteToken()
}
