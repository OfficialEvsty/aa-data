package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ITenantChainRepository interface {
	Add(context.Context, domain.TenantChain) error
	Remove(context.Context, uuid.UUID) error
	GetAllByTenant(context.Context, uuid.UUID) ([]*domain.TenantChain, error)
	CheckTenantAttachment(context.Context, uuid.UUID, []uuid.UUID) error
	WithTx(tx *sql.Tx) ITenantChainRepository
}
