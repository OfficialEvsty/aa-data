package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/interface/junction"
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
	query := `INSERT INTO aa_guild_nicknames (guild_id, nickname_id) VALUES ($1, $2) ON CONFLICT (guild_id, nickname_id) DO NOTHING;`
	_, err := r.exec.ExecContext(ctx, query, guildID, nicknameID)
	if err != nil {
		return err
	}
	return nil
}

func (r *GuildNicknameRepository) GetGuild(context.Context, *domain.AANickname) (*domain.GuildNickname, error) {
	return nil, nil
}
func (r *GuildNicknameRepository) GetMembers(ctx context.Context, guildID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT nickname_id FROM aa_guild_nicknames WHERE guild_id = $1;`
	rows, err := r.exec.QueryContext(ctx, query, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var guildNicknameList []uuid.UUID
	for rows.Next() {
		var nicknameID uuid.UUID
		if err = rows.Scan(&nicknameID); err != nil {
			return nil, err
		}
		guildNicknameList = append(guildNicknameList, nicknameID)
	}
	return guildNicknameList, nil
}
func (r *GuildNicknameRepository) ExcludeMember(context.Context, *domain.AANickname) error {
	return nil
}

func (r *GuildNicknameRepository) WithTx(tx *sql.Tx) junction_repos2.IGuildNicknameRepository {
	return &GuildNicknameRepository{tx}
}
