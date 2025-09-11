package domain

import (
	"github.com/google/uuid"
	"time"
)

type NicknameChain struct {
	ChainID       uuid.UUID  `json:"chain_id"`
	ParentChainID *uuid.UUID `json:"parent_chain_id"`
	NicknameID    uuid.UUID  `json:"nickname_id"`
	ChainedAt     time.Time  `json:"chained_at"`
	Active        bool       `json:"active"`
}
