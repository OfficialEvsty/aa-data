package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
)

// IItemRepository provides crud operations under items table
type IItemRepository interface {
	Add(context.Context, domain.AAItemTemplate) (*domain.AAItemTemplate, error)
	Remove(context.Context, int64) error
	GetByID(context.Context, int64) (*domain.AAItemTemplate, error)
	List(context.Context) ([]*domain.AAItemTemplate, error)
	WithTx(*sql.Tx) IItemRepository
}
