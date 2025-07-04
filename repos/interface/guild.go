package repos

import (
	"context"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// IGuildRepository provides operations with guild in db
type IGuildRepository interface {
	GetByName(ctx context.Context, serverID uuid.UUID, name string) (*domain.AAGuild, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AAGuild, error)
	List(ctx context.Context) ([]*domain.AAGuild, error)
	Add(ctx context.Context, guild domain.AAGuild) (*domain.AAGuild, error)
}
