package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ITenantRepository interface {
	Add(context.Context, domain.Tenant) (*domain.Tenant, error)
	Remove(context.Context, uuid.UUID) error
	Update(context.Context, domain.Tenant) (*domain.Tenant, error)
	GetByID(context.Context, uuid.UUID) (*domain.Tenant, error)
	GetByOwnerID(context.Context, uuid.UUID) (*domain.Tenant, error)
	GetOwnerID(context.Context, uuid.UUID) (uuid.UUID, error)
	All(context.Context) ([]*domain.Tenant, error)
	WithTx(*sql.Tx) ITenantRepository
}
