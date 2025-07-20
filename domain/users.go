package domain

import (
	"github.com/google/uuid"
	"time"
)

// Tenant contains necessary aa_guilds as one ally unit
type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   uuid.UUID `json:"owner_id"`
}

type TenantGuild struct {
	TenantID uuid.UUID `json:"tenant_id"`
	GuildID  uuid.UUID `json:"guild_id"`
	JoinedAt time.Time `json:"joined_at"`
}

// User record
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

// UserNickname record bounds user and their related nicknames
type UserNickname struct {
	UserID     uuid.UUID `json:"user_id"`
	NicknameID uuid.UUID `json:"nickname_id"`
}
