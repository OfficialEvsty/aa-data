package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
	"time"
)

type ILunarkRepository interface {
	Add(context.Context, domain.Lunark) error
	Update(context.Context, domain.Lunark) error
	UpdateEndDate(context.Context, uuid.UUID, time.Time) error
	Close(context.Context, uuid.UUID) error
	GetByID(context.Context, uuid.UUID) (*domain.Lunark, error)
	WithTx(*sql.Tx) ILunarkRepository
}
