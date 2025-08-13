package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
)

type IRaidItemRepository interface {
	AddOrUpdateItems(context.Context, uuid.UUID, []*serializable.DropItem) error
	RemoveItems(context.Context, uuid.UUID, []int) error
	RemoveItemsByRaidID(context.Context, uuid.UUID) error
	GetItems(context.Context, uuid.UUID) ([]*serializable.DropItem, error)
	WithTx(tx *sql.Tx) IRaidItemRepository
}
