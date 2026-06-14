package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/encryption"
	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
)

func (s *serv) RegisterUser(ctx context.Context, name, email, password, role string) error {
	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return ErrUserAlreadyExists
	}

	hashed, err := encryption.Hash(password)
	if err != nil {
		return err
	}

	return s.repo.SaveUser(ctx, name, email, hashed, role)
}

func (s *serv) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := encryption.CheckPassword(u.PasswordHash, password); err != nil {
		return nil, ErrInvalidPassword
	}

	return &models.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &models.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
