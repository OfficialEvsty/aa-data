package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

type IRaidRepository interface {
	Add(context.Context, domain.Raid) error
	Update(context.Context, domain.Raid) error
	UpdateTiming(context.Context, uuid.UUID, time.Time) error
	UpdateAttendance(context.Context, uuid.UUID, int) error
	UpdateStatus(context.Context, uuid.UUID, serializable.Status) error
	UpdateEndDateAndStatus(
		context.Context,
		uuid.UUID,
		time.Time,
		serializable.Status,
	) error
	Remove(context.Context, uuid.UUID) error
	GetById(context.Context, uuid.UUID) (*domain.Raid, error)
	WithTx(*sql.Tx) IRaidRepository
}
