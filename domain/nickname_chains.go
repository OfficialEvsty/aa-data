package domain

import (
	"github.com/google/uuid"
	"time"
)

type TenantChain struct {
	TenantID uuid.UUID `json:"tenant_id"`
	ChainID  uuid.UUID `json:"chain_id"`
}

type NicknameChain struct {
	ChainID       uuid.UUID  `json:"chain_id"`
	ParentChainID *uuid.UUID `json:"parent_chain_id"`
	NicknameID    uuid.UUID  `json:"nickname_id"`
	CreatedAt     time.Time  `json:"created_at"`
	ChainedAt     *time.Time `json:"chained_at"`
	Active        bool       `json:"active"`
}
