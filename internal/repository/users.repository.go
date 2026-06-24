package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveUser(ctx context.Context, name, email, passwordHash, role string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4)`,
		name, email, passwordHash, role,
	)
	return err
}

func buildUserWhere(nameQuery string, activeFilter *bool) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if nameQuery != "" {
		args = append(args, "%"+nameQuery+"%")
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d)", len(args), len(args)))
	}
	if activeFilter != nil {
		args = append(args, *activeFilter)
		conditions = append(conditions, fmt.Sprintf("active = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *repo) CountUsers(ctx context.Context, nameQuery string, activeFilter *bool) (int, error) {
	where, args := buildUserWhere(nameQuery, activeFilter)
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM users`+where, args...)
	return count, err
}

func (r *repo) GetUsers(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]entity.User, error) {
	where, args := buildUserWhere(nameQuery, activeFilter)
	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT * FROM users%s ORDER BY name LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	var users []entity.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	return users, err
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	var u entity.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) UpdateUserActive(ctx context.Context, id string, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET active=$1, updated_at=NOW() WHERE id=$2`,
		active, id,
	)
	return err
}
