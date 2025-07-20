package domain

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
