package repos

import (
	"context"
	"database/sql"
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
	guild.ID = uuid.New()
	query := `INSERT INTO aa_guilds (id, name, server_id) 
              VALUES ($1, $2, $3) ON CONFLICT (server_id, name) DO UPDATE SET name = EXCLUDED.name RETURNING id, name, server_id;`
	res := r.exec.QueryRowContext(ctx, query, guild.ID, guild.Name, guild.ServerID)
	err := res.Scan(&result.ID, &result.Name, &result.ServerID)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var guild domain.AAGuild
		err = rows.Scan(&guild.ID, &guild.Name, &guild.ServerID)
		if err != nil {
			return nil, err
		}
		result = append(result, &guild)
	}
	return result, nil
}

func (r *GuildRepository) WithTx(tx *sql.Tx) repos.IGuildRepository {
	return &GuildRepository{tx}
}
