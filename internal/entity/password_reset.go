package entity

import (
	"database/sql"
	"time"
)

type PasswordReset struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	Token     string       `db:"token"`
	ExpiresAt time.Time    `db:"expires_at"`
	UsedAt    sql.NullTime `db:"used_at"`
	CreatedAt time.Time    `db:"created_at"`
}
