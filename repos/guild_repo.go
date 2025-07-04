package repos

import (
	"context"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
)

// GuildRepository guild's repository implementation
type GuildRepository struct {
	exec db.ISqlExecutor
}

// NewGuildRepository creates instance of GuildRepository
func NewGuildRepository(executor db.ISqlExecutor) *GuildRepository {
	return &GuildRepository{
		exec: executor,
	}
}

// Add saves domain.AAGuild in table aa_guilds
func (r *GuildRepository) Add(ctx context.Context, guild domain.AAGuild) (*domain.AAGuild, error) {
	var result domain.AAGuild
	query := `INSERT INTO aa_guilds (id, name, server_id) 
              VALUES ($1, $2, $3)`
	res, err := r.exec.QueryContext(ctx, query, guild.ID, guild.Name, guild.ServerID)
	if err != nil {
		return nil, fmt.Errorf("error while inserting guild: %v", err)
	}
	return &result, nil
}
