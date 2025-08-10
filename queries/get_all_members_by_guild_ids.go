package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetAllMembersByGuildIDs struct {
	exec db.ISqlExecutor
}

func NewGetAllMembersByGuildIDs(exec db.ISqlExecutor) *GetAllMembersByGuildIDs {
	return &GetAllMembersByGuildIDs{exec: exec}
}

func (q *GetAllMembersByGuildIDs) Handle(ctx context.Context, guildIDs []uuid.UUID) ([]*usecase.GuildMemberDTO, error) {
	query := `SELECT g.id, g.name, n.id, n.name
              FROM aa_guild_nicknames gn
              JOIN aa_guilds g ON g.id = gn.guild_id
              JOIN aa_nicknames n ON n.id = gn.nickname_id
              WHERE g.id = ANY($1)`
	rows, err := q.exec.QueryContext(ctx, query, guildIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	members := make([]*usecase.GuildMemberDTO, 0)
	for rows.Next() {
		var member usecase.GuildMemberDTO
		err = rows.Scan(
			&member.GuildID,
			&member.GuildName,
			&member.NicknameID,
			&member.Name,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}
