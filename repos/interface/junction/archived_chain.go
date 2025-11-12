package junction_repos

import (
	"context"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type IArchivedChainRepository interface {
	Add(ctx context.Context, aChain domain.ArchivedChain) error
	AddMany(ctx context.Context, aChains []domain.ArchivedChain) error
	List(ctx context.Context, tenantID uuid.UUID) ([]*domain.ArchivedChain, error)
	IsArchived(ctx context.Context, chainID uuid.UUID) (bool, error)
	Remove(ctx context.Context, chainID uuid.UUID) error
	RemoveMany(ctx context.Context, chainIDs []uuid.UUID) error
}
