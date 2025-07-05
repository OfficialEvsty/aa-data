package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

// GuildNicknameRepository provides links between guild and nickname
type GuildNicknameRepository struct {
	exec db.ISqlExecutor
}

func NewGuildNicknameRepository(exec db.ISqlExecutor) *GuildNicknameRepository {
	return &GuildNicknameRepository{exec}
}

// Add links guild and nickname
func (r *GuildNicknameRepository) Add(ctx context.Context, guildID uuid.UUID, nicknameID uuid.UUID) error {
	query := `INSERT INTO aa_guild_nicknames (guild_id, nickname_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	_, err := r.exec.ExecContext(ctx, query, guildID, nicknameID)
	if err != nil {
		return err
	}
	return nil
}

func (r *GuildNicknameRepository) GetGuild(context.Context, *domain.AANickname) (*domain.GuildNickname, error) {
	return nil, nil
}
func (r *GuildNicknameRepository) GetMembers(context.Context, *domain.AAGuild) ([]*domain.GuildNickname, error) {
	return nil, nil
}
func (r *GuildNicknameRepository) ExcludeMember(context.Context, *domain.AANickname) error {
	return nil
}

func (r *GuildNicknameRepository) WithTx(tx *sql.Tx) repos.IGuildNicknameRepository {
	return &GuildNicknameRepository{tx}
}
