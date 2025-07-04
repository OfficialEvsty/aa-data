package queries

import (
	"context"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// GuildNicknamesQuery represents nicknames by specified guild's ids
type GuildNicknamesQuery struct {
	exec db.ISqlExecutor
}

func NewGuildNicknamesQuery(exec db.ISqlExecutor) *GuildNicknamesQuery {
	return &GuildNicknamesQuery{exec}
}

func (q *GuildNicknamesQuery) GetNicknamesByGuildID(ctx context.Context, guildID uuid.UUID) ([]domain.GuildNickname, error) {
	query := `SELECT g.id, g.name, g.server_id, n.id, n.name, n.server_id, n.created_at, gn.member_at 
			  FROM aa_guild_nicknames AS gn 
			  JOIN aa_guilds AS g ON gn.guild_id = g.id
			  JOIN aa_nicknames AS n ON gn.nickname_id = n.id
			  WHERE gn.guild_id = $1`
	rows, err := q.exec.QueryContext(ctx, query, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.GuildNickname
	for rows.Next() {
		var gn domain.GuildNickname
		if err = rows.Scan(&gn.Guild.ID, &gn.Guild.Name, &gn.Guild.ServerID,
			&gn.Nickname.ID, &gn.Nickname.Name, &gn.Nickname.ServerID,
			&gn.Nickname.CreatedAt, &gn.MemberAt); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		result = append(result, gn)
	}
	return result, nil
}

func (q *GuildNicknamesQuery) GetNicknamesByGuildIDs(ctx context.Context, guildIDs []uuid.UUID) (map[uuid.UUID][]domain.GuildNickname, error) {
	query := `SELECT g.id, g.name, g.server_id, n.id, n.name, n.server_id, n.created_at, gn.member_at 
			  FROM aa_guild_nicknames AS gn 
			  JOIN aa_guilds AS g ON gn.guild_id = g.id
			  JOIN aa_nicknames AS n ON gn.nickname_id = n.id
			  WHERE gn.guild_id = ANY($1)`
	rows, err := q.exec.QueryContext(ctx, query, guildIDs)
	if err != nil {
		return nil, fmt.Errorf("error while getting guild members: %v", err)
	}
	defer rows.Close()
	result := make(map[uuid.UUID][]domain.GuildNickname)
	for rows.Next() {
		var gn domain.GuildNickname
		if err = rows.Scan(&gn.Guild.ID, &gn.Guild.Name, &gn.Guild.ServerID,
			&gn.Nickname.ID, &gn.Nickname.Name, &gn.Nickname.ServerID,
			&gn.Nickname.CreatedAt, &gn.MemberAt); err != nil {
			return nil, fmt.Errorf("error while scan guild members: %v", err)
		}
		result[gn.Guild.ID] = append(result[gn.Guild.ID], gn)
	}
	return result, nil
}
