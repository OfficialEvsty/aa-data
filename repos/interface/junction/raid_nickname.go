package junction_repos

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type IRaidNicknameRepository interface {
	AddNicknames(context.Context, uuid.UUID, []uuid.UUID) error
	RemoveNicknames(context.Context, uuid.UUID, []uuid.UUID) error
	ClearNicknamesByRaidID(context.Context, uuid.UUID) error
	WithTx(*sql.Tx) IRaidNicknameRepository
}
