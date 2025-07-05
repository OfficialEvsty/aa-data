package repos

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
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
	err = res.Scan(&result.ID, &result.Name, &result.ServerID)
	if err != nil {
		return nil, fmt.Errorf("error while scanning guild: %v", err)
	}
	return &result, nil
}

func (r *GuildRepository) GetByName(ctx context.Context, serverID uuid.UUID, name string) (*domain.AAGuild, error) {
	var result domain.AAGuild
	return &result, nil
}
func (r *GuildRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AAGuild, error) {
	var result domain.AAGuild
	return &result, nil
}
func (r *GuildRepository) List(ctx context.Context) ([]*domain.AAGuild, error) {
	var result []*domain.AAGuild
	rows, err := r.exec.QueryContext(ctx, "SELECT id, name, server_id FROM aa_guilds")
	if err != nil {
		return nil, fmt.Errorf("error while listing guilds: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var guild domain.AAGuild
		err = rows.Scan(&guild.ID, &guild.Name, &guild.ServerID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning guild: %v", err)
		}
		result = append(result, &guild)
	}
	return result, nil
}

func (r *GuildRepository) WithTx(tx *sql.Tx) repos.IGuildRepository {
	return &GuildRepository{tx}
}
