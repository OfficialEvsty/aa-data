package usecase

import "github.com/google/uuid"

type GuildMemberDTO struct {
	GuildID    uuid.UUID `json:"guild_id"`
	GuildName  string    `json:"guild_name"`
	NicknameID uuid.UUID `json:"nickname_id"`
	Name       string    `json:"name"`
}
