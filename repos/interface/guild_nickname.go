package repos

import (
	"context"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// IGuildNicknameRepository provides operation with guild members in db
type IGuildNicknameRepository interface {
	GetGuild(context.Context, *domain.AANickname) (*domain.GuildNickname, error)
	GetMembers(context.Context, *domain.AAGuild) ([]*domain.GuildNickname, error)
	Add(context.Context, uuid.UUID, uuid.UUID) error
	ExcludeMember(context.Context, *domain.AANickname) error
}
