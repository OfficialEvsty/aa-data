package repos

import (
	"context"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// INicknameRepository interface for aa-nicknames interaction in db
type INicknameRepository interface {
	GetByName(ctx context.Context, serverID uuid.UUID, name string) (*domain.AANickname, error)
	Create(ctx context.Context, nick domain.AANickname) (*domain.AANickname, error)
}
