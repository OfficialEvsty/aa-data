package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type IPaymentRepository interface {
	AddMany(context.Context, []domain.Payment) error
	Clear(context.Context, uuid.UUID) error
	GetAllBySalaryID(context.Context, uuid.UUID) ([]*domain.Payment, error)
	GetAllByChainID(context.Context, uuid.UUID) ([]*domain.Payment, error)
	WithTx(*sql.Tx) IPaymentRepository
}
