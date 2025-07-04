package queries

import (
	"context"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
	"time"
)

type ServersGuildsNicknamesQuery struct {
	exec db.ISqlExecutor
}

func NewServersGuildsNicknamesQuery(exec db.ISqlExecutor) *ServersGuildsNicknamesQuery {
	return &ServersGuildsNicknamesQuery{exec}
}

// GetServerData provides data about guilds and nicknames on server
func (q *ServersGuildsNicknamesQuery) GetServerData(ctx context.Context) ([]*usecase.ServerData, error) {
	var servers []*usecase.ServerData
	query := `SELECT s.id, s.name, s.external_id,
			  g.id, g.name,
			  n.id, n.name, n.created_at, gn.member_at
 			  FROM aa_servers s 
 			  JOIN aa_guilds g ON s.id = g.server_id
 			  JOIN aa_guild_nicknames AS gn ON gn.guild_id = g.id
 			  JOIN aa_nicknames AS n ON n.id = gn.nickname_id`
	rows, err := q.exec.QueryContext(ctx, query)
	if err != nil {
		return servers, fmt.Errorf("error while getting server data: %v", err)
	}
	defer rows.Close()
	var (
		serverMap = make(map[uuid.UUID]*usecase.ServerData)
	)
	for rows.Next() {
		var (
			serverID         uuid.UUID
			serverName       string
			serverExternalID string

			guildID   uuid.UUID
			guildName string

			nicknameID   uuid.UUID
			nicknameName string
			createdAt    time.Time

			memberAt time.Time
		)
		if err = rows.Scan(
			&serverID, &serverName, &serverExternalID,
			&guildID, &guildName,
			&nicknameID, &nicknameName, &createdAt,
			&memberAt,
		); err != nil {
			return nil, err
		}
		serverData, ok := serverMap[serverID]
		if !ok {
			serverData = &usecase.ServerData{
				Server: &domain.AAServer{
					ID:         serverID,
					Name:       serverName,
					ExternalID: serverExternalID,
				},
				Guilds: []*usecase.GuildData{},
			}
			serverMap[serverID] = serverData
		}

		// Guild map scoped per server
		var guildMap map[uuid.UUID]*usecase.GuildData
		if serverData.Guilds == nil {
			guildMap = make(map[uuid.UUID]*usecase.GuildData)
		} else {
			guildMap = make(map[uuid.UUID]*usecase.GuildData)
			for _, g := range serverData.Guilds {
				guildMap[g.Guild.ID] = g
			}
		}

		// Get or create GuildData
		guildData, ok := guildMap[guildID]
		if !ok {
			guildData = &usecase.GuildData{
				Guild: &domain.AAGuild{
					ID:       guildID,
					Name:     guildName,
					ServerID: serverID,
				},
				Nicknames: []*domain.GuildNickname{},
			}
			serverData.Guilds = append(serverData.Guilds, guildData)
			guildMap[guildID] = guildData
		}

		// Add nickname
		guildData.Nicknames = append(guildData.Nicknames, &domain.GuildNickname{
			Nickname: domain.AANickname{
				ID:        nicknameID,
				Name:      nicknameName,
				CreatedAt: createdAt,
				ServerID:  serverID,
			},
			MemberAt: memberAt,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// ✨ Результат — список всех серверов
	var result []*usecase.ServerData
	for _, s := range serverMap {
		result = append(result, s)
	}
	return result, nil
}
