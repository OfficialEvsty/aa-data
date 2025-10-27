package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type IRequestRepository interface {
	Add(context.Context, domain.Request) error
	Remove(context.Context, uuid.UUID) error
	Accept(context.Context, uuid.UUID, uuid.UUID) error
	Decline(context.Context, uuid.UUID, uuid.UUID) error
	Get(context.Context, uuid.UUID) (*domain.Request, error)
	WithTx(tx *sql.Tx) IRequestRepository
}
