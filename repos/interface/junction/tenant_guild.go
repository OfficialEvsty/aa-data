package junction_repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type ITenantGuildRepository interface {
	AddGuild(context.Context, domain.TenantGuild) (*domain.TenantGuild, error)
	RemoveGuild(context.Context, uuid.UUID) error
	AddGuilds(context.Context, uuid.UUID, []uuid.UUID) error
	RemoveGuilds(context.Context, uuid.UUID, []uuid.UUID) error
	All(context.Context, uuid.UUID) ([]*domain.TenantGuild, error)
	GetByGuildID(context.Context, uuid.UUID) (*domain.TenantGuild, error)
	WithTx(*sql.Tx) ITenantGuildRepository
}
