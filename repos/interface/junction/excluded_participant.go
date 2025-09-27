package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type IExcludedParticipantRepository interface {
	Add(context.Context, usecase.ExcludedParticipant) error
	Remove(context.Context, uuid.UUID, uuid.UUID) error
	All(context.Context, uuid.UUID) ([]*usecase.ExcludedParticipant, error)
	WithTx(*sql.Tx) IExcludedParticipantRepository
}
