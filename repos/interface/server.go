package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type IServerRepository interface {
	Add(context.Context, domain.AAServer) (*domain.AAServer, error)
	// GetByExternalID for working with official RU Archeage server's ids
	GetByExternalID(context.Context, string) (*domain.AAServer, error)
	List(context.Context) ([]*domain.AAServer, error)
	Remove(context.Context, uuid.UUID) error
	WithTx(*sql.Tx) IServerRepository
}
