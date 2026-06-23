package models

type UserInvite struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	InviteURL string `json:"invite_url"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}
