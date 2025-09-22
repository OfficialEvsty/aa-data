package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type IRaidEventRepository interface {
	Add(context.Context, domain.RaidEvent) error
	AddMany(context.Context, uuid.UUID, []int) error
	Remove(context.Context, uuid.UUID, int) error
	All(context.Context, uuid.UUID) ([]*domain.RaidEvent, error)
	AllByRaidIDs(context.Context, []uuid.UUID) (usecase.RaidEvents, error)
	WithTx(*sql.Tx) IRaidEventRepository
}
