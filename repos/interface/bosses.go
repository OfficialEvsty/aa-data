package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
)

// IBossesRepository crud operations under bosses and their drop
type IBossesRepository interface {
	Add(context.Context, domain.AABoss) (*domain.AABoss, error)
	Remove(context.Context, int64) error
	GetByID(context.Context, int64) (*domain.AABoss, error)
	List(context.Context) ([]*domain.AABoss, error)
	WithTx(*sql.Tx) IBossesRepository
}
