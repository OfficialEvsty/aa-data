package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ILunarkRaidRepository interface {
	Add(context.Context, domain.LunarkRaid) error
	Remove(context.Context, uuid.UUID) error
	All(context.Context, uuid.UUID) ([]*domain.LunarkRaid, error)
	GetByID(context.Context, uuid.UUID) (*domain.LunarkRaid, error)
	WithTx(*sql.Tx) ILunarkRaidRepository
}
