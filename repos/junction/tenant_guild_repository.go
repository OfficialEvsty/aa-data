package junction_repos

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
)

type TenantGuildRepository struct {
	exec db.ISqlExecutor
}

func NewTenantGuildRepository(exec db.ISqlExecutor) *TenantGuildRepository {
	return &TenantGuildRepository{exec}
}

func (r *TenantGuildRepository) AddGuilds(ctx context.Context, tenantID uuid.UUID, guildIDs []uuid.UUID) error {
	valueStrings := make([]string, 0, len(guildIDs))
	valueArgs := make([]interface{}, 0, len(guildIDs)*2)

	for i, id := range guildIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, tenantID, id)
	}

	stmt := fmt.Sprintf("INSERT INTO aa_tenant_guilds (tenant_id, guild_id) VALUES %s", strings.Join(valueStrings, ","))
	_, err := r.exec.ExecContext(ctx, stmt, valueArgs...)
	return err
}
func (r *TenantGuildRepository) RemoveGuilds(ctx context.Context, tenantID uuid.UUID, guildIDs []uuid.UUID) error {
	query := `DELETE FROM aa_tenant_guilds WHERE tenant_id = $1 AND guild_id = ANY($2)`
	_, err := r.exec.ExecContext(ctx, query, tenantID, pq.Array(guildIDs))
	if err != nil {
		return err
	}
	return nil
}

func (r *TenantGuildRepository) AddGuild(ctx context.Context, allyGuild domain.TenantGuild) (*domain.TenantGuild, error) {
	query := `INSERT INTO aa_tenant_guilds (tenant_id, guild_id)
			  VALUES ($1, $2) RETURNING tenant_id, guild_id, joined_at;`
	err := r.exec.QueryRowContext(ctx, query, allyGuild.TenantID, allyGuild.GuildID).Scan(&allyGuild.TenantID, &allyGuild.GuildID, &allyGuild.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &allyGuild, nil
}
func (r *TenantGuildRepository) RemoveGuild(ctx context.Context, guildID uuid.UUID) error {
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
