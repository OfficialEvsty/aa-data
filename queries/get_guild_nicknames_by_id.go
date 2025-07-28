package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-shared/golinq"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type GuildNicknameDTO struct {
	ServerID   uuid.UUID `json:"server_id"`
	GuildID    uuid.UUID `json:"guild_id"`
	GuildName  string    `json:"guild_name"`
	NicknameID uuid.UUID `json:"nickname_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetGuildNicknamesByIdQuery struct {
	exec db.ISqlExecutor
}

func NewGetGuildNicknamesByIdQuery(exec db.ISqlExecutor) GetGuildNicknamesByIdQuery {
	return GetGuildNicknamesByIdQuery{
		exec: exec,
	}
}

func (q GetGuildNicknamesByIdQuery) Handle(ctx context.Context, guildIDs []uuid.UUID) ([]*GuildNicknameDTO, error) {

	guildIDsStr := golinq.Map(guildIDs, func(id uuid.UUID) string {
		return id.String()
	})
	query := `SELECT n.server_id, gn.guild_id, g.name, gn.nickname_id, n.name, n.created_at
			  FROM aa_guild_nicknames AS gn
			  JOIN aa_guilds AS g ON gn.guild_id = g.id
			  JOIN aa_nicknames AS n ON gn.nickname_id = n.id
			  WHERE gn.guild_id = ANY($1::uuid[])`
	rows, err := q.exec.QueryContext(ctx, query, pq.Array(guildIDsStr))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var guildNicknames []*GuildNicknameDTO
	for rows.Next() {
		var guildNicknameDTO GuildNicknameDTO
		err = rows.Scan(
			&guildNicknameDTO.ServerID,
			&guildNicknameDTO.GuildID,
			&guildNicknameDTO.GuildName,
			&guildNicknameDTO.NicknameID,
			&guildNicknameDTO.Name,
			&guildNicknameDTO.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		guildNicknames = append(guildNicknames, &guildNicknameDTO)
	}
	return guildNicknames, nil
}
