package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// IFinishedPublish when publish finishes results stored
type IFinishedPublish interface {
	Add(context.Context, domain.FinishedPublish) (*domain.FinishedPublish, error)
	Remove(context.Context, uuid.UUID) error
	Get(context.Context, uuid.UUID) (*domain.FinishedPublish, error)
	WithTx(*sql.Tx) IFinishedPublish
}
