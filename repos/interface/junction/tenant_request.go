package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ITenantRequestRepository interface {
	Add(context.Context, domain.TenantRequest) error
	Remove(context.Context, uuid.UUID) error
	GetAllByTenantID(context.Context, uuid.UUID) ([]uuid.UUID, error)
	WithTx(*sql.Tx) ITenantRequestRepository
}
