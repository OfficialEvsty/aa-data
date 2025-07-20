package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type TenantGuildRepository struct {
	exec db.ISqlExecutor
}

func NewTenantGuildRepository(exec db.ISqlExecutor) *TenantGuildRepository {
	return &TenantGuildRepository{exec}
}

func (r *TenantGuildRepository) Add(ctx context.Context, allyGuild domain.TenantGuild) (*domain.TenantGuild, error) {
	query := `INSERT INTO aa_tenant_guilds (tenant_id, guild_id)
			  VALUES ($1, $2) RETURNING tenant_id, guild_id, joined_at;`
	err := r.exec.QueryRowContext(ctx, query, allyGuild.TenantID, allyGuild.GuildID).Scan(&allyGuild.TenantID, &allyGuild.GuildID, &allyGuild.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &allyGuild, nil
}
func (r *TenantGuildRepository) Remove(ctx context.Context, guildID uuid.UUID) error {
	query := `DELETE FROM aa_tenant_guilds WHERE guild_id = $1;`
	_, err := r.exec.ExecContext(ctx, query, guildID)
	if err != nil {
		return err
	}
	return nil
}
func (r *TenantGuildRepository) All(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantGuild, error) {
	query := `SELECT tenant_id, guild_id, joined_at FROM aa_tenant_guilds WHERE tenant_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tenantGuilds []*domain.TenantGuild
	for rows.Next() {
		var tenantGuild domain.TenantGuild
		err = rows.Scan(&tenantGuild.TenantID, &tenantGuild.GuildID, &tenantGuild.JoinedAt)
		if err != nil {
			return nil, err
		}
		tenantGuilds = append(tenantGuilds, &tenantGuild)
	}
	return tenantGuilds, nil
}
func (r *TenantGuildRepository) GetByGuildID(ctx context.Context, guildID uuid.UUID) (*domain.TenantGuild, error) {
	var tenantGuild domain.TenantGuild
	query := `SELECT tenant_id, guild_id, joined_at FROM aa_tenant_guilds WHERE guild_id = $1`
	row := r.exec.QueryRowContext(ctx, query, guildID)
	err := row.Scan(&tenantGuild.TenantID, &tenantGuild.GuildID, &tenantGuild.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &tenantGuild, nil
}
func (r *TenantGuildRepository) WithTx(tx *sql.Tx) junction_repos.ITenantGuildRepository {
	return &TenantGuildRepository{tx}
}
