package entity

import (
	"database/sql"
	"time"
)

type UserInvite struct {
	ID         string         `db:"id"`
	Email      string         `db:"email"`
	Role       string         `db:"role"`
	Token      string         `db:"token"`
	InvitedBy  sql.NullString `db:"invited_by"`
	ExpiresAt  time.Time      `db:"expires_at"`
	AcceptedAt sql.NullTime   `db:"accepted_at"`
	CreatedAt  time.Time      `db:"created_at"`
}
