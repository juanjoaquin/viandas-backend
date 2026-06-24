package service

import (
	"context"
	"errors"
	"time"

	"github.com/juanjoaquin/viandas-backend/encryption"
	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

func (s *serv) CreateRefreshToken(ctx context.Context, userID string) (string, error) {
	tokenStr, err := encryption.SignedRefreshToken(userID)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(encryption.RefreshTokenDuration)
	if err := s.repo.SaveRefreshToken(ctx, userID, tokenStr, expiresAt); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *serv) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := encryption.ParseRefreshJWT(refreshToken)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	stored, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil || stored == nil {
		return "", "", ErrInvalidRefreshToken
	}

	if time.Now().After(stored.ExpiresAt) {
		_ = s.repo.DeleteRefreshToken(ctx, refreshToken)
		return "", "", ErrInvalidRefreshToken
	}

	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	if !user.Active {
		_ = s.repo.DeleteRefreshToken(ctx, refreshToken)
		return "", "", ErrUserInactive
	}

	accessToken, err := encryption.SignedLoginToken(&models.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
	if err != nil {
		return "", "", err
	}

	// Rotación: invalidar el refresh token usado y emitir uno nuevo
	_ = s.repo.DeleteRefreshToken(ctx, refreshToken)

	newRefreshToken, err := s.CreateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *serv) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, refreshToken)
}
