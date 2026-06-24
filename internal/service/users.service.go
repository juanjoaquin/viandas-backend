package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/juanjoaquin/viandas-backend/encryption"
	"github.com/juanjoaquin/viandas-backend/internal/entity"
	"github.com/juanjoaquin/viandas-backend/internal/models"
	"github.com/juanjoaquin/viandas-backend/internal/roles"
)

var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserInactive             = errors.New("account disabled")
	ErrUserCannotDeactivateSelf = errors.New("cannot deactivate your own account")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrInvalidInvite            = errors.New("invalid invite")
	ErrInviteExpired            = errors.New("invite expired")
	ErrInviteAlreadyAccepted    = errors.New("invite already accepted")
	ErrInviteRoleNotSupported   = errors.New("invite role not supported")
)

func (s *serv) RegisterUser(ctx context.Context, name, email, password, role string) error {
	email = normalizeEmail(email)
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

func (s *serv) InviteUser(ctx context.Context, email, role, invitedBy string) (*models.UserInvite, error) {
	email = normalizeEmail(email)
	if role != roles.Employee {
		return nil, ErrInviteRoleNotSupported
	}

	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	token, err := generateInviteToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(72 * time.Hour)
	invite, err := s.repo.SaveUserInvite(ctx, email, role, token, invitedBy, expiresAt)
	if err != nil {
		return nil, err
	}

	inviteURL := s.inviteURL(token)
	if err := s.mailer.SendInvite(ctx, email, inviteURL); err != nil {
		return nil, err
	}

	return inviteToModel(invite, inviteURL), nil
}

func (s *serv) RegisterWithInvite(ctx context.Context, token, name, password string) error {
	invite, err := s.repo.GetUserInviteByToken(ctx, strings.TrimSpace(token))
	if err != nil {
		return ErrInvalidInvite
	}
	if invite.AcceptedAt.Valid {
		return ErrInviteAlreadyAccepted
	}
	if time.Now().After(invite.ExpiresAt) {
		return ErrInviteExpired
	}

	existing, _ := s.repo.GetUserByEmail(ctx, invite.Email)
	if existing != nil {
		return ErrUserAlreadyExists
	}

	hashed, err := encryption.Hash(password)
	if err != nil {
		return err
	}

	if err := s.repo.SaveUser(ctx, name, invite.Email, hashed, invite.Role); err != nil {
		return err
	}

	return s.repo.AcceptUserInvite(ctx, invite.ID)
}

func (s *serv) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	u, err := s.repo.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := encryption.CheckPassword(u.PasswordHash, password); err != nil {
		return nil, ErrInvalidPassword
	}

	if !u.Active {
		return nil, ErrUserInactive
	}

	return userEntityToModel(u), nil
}

func userEntityToModel(u *entity.User) *models.User {
	return &models.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (s *serv) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return userEntityToModel(u), nil
}

func (s *serv) CountUsers(ctx context.Context, nameQuery string, activeFilter *bool) (int, error) {
	return s.repo.CountUsers(ctx, nameQuery, activeFilter)
}

func (s *serv) GetUsers(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]models.User, error) {
	entities, err := s.repo.GetUsers(ctx, nameQuery, activeFilter, offset, limit)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(entities))
	for i, u := range entities {
		users[i] = *userEntityToModel(&u)
	}
	return users, nil
}

func (s *serv) UpdateUserActive(ctx context.Context, id string, active bool, requestingUserID string) error {
	if !active && id == requestingUserID {
		return ErrUserCannotDeactivateSelf
	}

	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.repo.UpdateUserActive(ctx, id, active); err != nil {
		return err
	}

	if !active {
		_ = s.repo.DeleteRefreshTokensByUserID(ctx, id)
	}

	return nil
}

func (s *serv) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := s.repo.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil, ErrUserNotFound
	}

	return userEntityToModel(u), nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func generateInviteToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *serv) inviteURL(token string) string {
	baseURL := strings.TrimRight(s.settings.FrontendURL, "/")
	return fmt.Sprintf("%s/register?token=%s", baseURL, url.QueryEscape(token))
}

func inviteToModel(invite *entity.UserInvite, inviteURL string) *models.UserInvite {
	return &models.UserInvite{
		ID:        invite.ID,
		Email:     invite.Email,
		Role:      invite.Role,
		InviteURL: inviteURL,
		ExpiresAt: invite.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		CreatedAt: invite.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
