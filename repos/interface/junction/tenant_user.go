package junction_repos

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type ITenantUserRepository interface {
	Add(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) error
	Remove(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) error
	GetTenant(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetUserIDs(ctx context.Context, tenantID uuid.UUID) ([]uuid.UUID, error)
	CheckUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (bool, error)
	WithTx(*sql.Tx) ITenantUserRepository
}
