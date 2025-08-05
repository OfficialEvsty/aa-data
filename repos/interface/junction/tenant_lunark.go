package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ITenantLunarkRepository interface {
	Add(context.Context, domain.Journal) error
	Remove(context.Context, uuid.UUID) error
	All(context.Context, uuid.UUID) ([]*domain.Journal, error)
	GetByID(context.Context, uuid.UUID) (*domain.Journal, error)
	WithTx(*sql.Tx) ITenantLunarkRepository
}
