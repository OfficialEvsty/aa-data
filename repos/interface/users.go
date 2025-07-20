package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// IUserRepository crud under users
type IUserRepository interface {
	AddOrUpdate(context.Context, domain.User) (*domain.User, error)
	GetByID(context.Context, uuid.UUID) (*domain.User, error)
	//Update(context.Context, domain.User) (*domain.User, error)
	Remove(context.Context, uuid.UUID) error
	WithTx(*sql.Tx) IUserRepository
}
