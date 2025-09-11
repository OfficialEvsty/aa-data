package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type IChainRepository interface {
	Add(context.Context, domain.NicknameChain) error
	GetChain(context.Context, uuid.UUID) ([]*domain.NicknameChain, error)
	GetRootChain(context.Context, uuid.UUID) (*domain.NicknameChain, error)
	GetActiveChainID(context.Context, uuid.UUID) (*domain.NicknameChain, error)
	GetNicknameChains(context.Context, uuid.UUID) ([]*domain.NicknameChain, error)
	AttachChain(ctx context.Context, parent uuid.UUID, child uuid.UUID) error
	DetachChain(context.Context, uuid.UUID) error
	UpdateStatus(context.Context, uuid.UUID, bool) error
	WithTx(*sql.Tx) IChainRepository
}
