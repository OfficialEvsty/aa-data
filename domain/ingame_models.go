package domain

import (
	"github.com/google/uuid"
	"time"
)

// AAServer in-game Archeage server
type AAServer struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ExternalID string    `json:"external_id"`
}

// AANickname in-game ArcheAge nickname
type AANickname struct {
	ID        uuid.UUID `json:"id"`
	ServerID  uuid.UUID `json:"server_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// AAGuild in-game guild with members
type AAGuild struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ServerID uuid.UUID `json:"server_id"`
}

// GuildNickname relation between nicknames and guilds
type GuildNickname struct {
	Nickname AANickname `json:"nickname"`
	Guild    AAGuild    `json:"guild"`
	MemberAt time.Time  `json:"member_at"`
}

type ServerNickname struct {
	Nickname AANickname `json:"nickname"`
	Server   AAServer   `json:"guild"`
}
