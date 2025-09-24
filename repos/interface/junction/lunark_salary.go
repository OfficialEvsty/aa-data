package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ILunarkSalaryRepository interface {
	Add(context.Context, domain.LunarkSalary) error
	Remove(context.Context, uuid.UUID) error
	GetSalaries(context.Context, uuid.UUID) ([]*domain.LunarkSalary, error)
	WithTx(tx *sql.Tx) ILunarkSalaryRepository
}
