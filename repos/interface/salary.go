package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ISalaryRepository interface {
	Add(context.Context, domain.Salary) error
	SafeDelete(context.Context, uuid.UUID) error
	Update(context.Context, domain.Salary) error
	GetByID(context.Context, uuid.UUID) (*domain.Salary, error)
	WithTx(*sql.Tx) ISalaryRepository
}
