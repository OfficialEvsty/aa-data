package usecase

import (
	"github.com/google/uuid"
	"time"
)

type ChainedNodeNickname struct {
	ParentChainID uuid.UUID `json:"parent_chain_id"`
	ChainID       uuid.UUID `json:"chain_id"`
	NicknameID    string    `json:"nickname_id"`
	OwnedAt       time.Time `json:"owned_at"`
	Active        bool      `json:"active"`
}

type ChainedParticipant struct {
	RootChainID       uuid.UUID             `json:"root_chain_id"`
	GuildID           uuid.UUID             `json:"guild_id"`
	GuildName         string                `json:"guild_name"`
	CurrentNicknameID uuid.UUID             `json:"nickname_id"`
	CurrentNickname   string                `json:"nickname"`
	History           []ChainedNodeNickname `json:"history"`
}

type AllChainedParticipants = map[uuid.UUID]*ChainedParticipant

// NicknameID keys refers to root chain id
type ChainedNicknames = map[uuid.UUID]uuid.UUID
