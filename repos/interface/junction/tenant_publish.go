package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ITenantPublishRepository interface {
	Add(ctx context.Context, publish domain.TenantPublish) (*domain.TenantPublish, error)
	Remove(ctx context.Context, publishID uuid.UUID) error
	GetByID(ctx context.Context, publishID uuid.UUID) (*domain.TenantPublish, error)
	All(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantPublish, error)
	RemoveAll(ctx context.Context, tenantID uuid.UUID) error
	WithTx(*sql.Tx) ITenantPublishRepository
}
