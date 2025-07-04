package repos

import (
	"context"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
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
		return fmt.Errorf("error linking guild and nickname: %v", err)
	}
	return nil
}
