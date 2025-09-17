package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type GetChainedParticipantsByTenantIDQuery struct {
	exec db.ISqlExecutor
}

func NewGetChainedParticipantsByTenantIDQuery(exec db.ISqlExecutor) *GetChainedParticipantsByTenantIDQuery {
	return &GetChainedParticipantsByTenantIDQuery{exec: exec}
}

func (q *GetChainedParticipantsByTenantIDQuery) Handle(ctx context.Context, activeChainIDs []uuid.UUID) (usecase.AllChainedParticipants, error) {
	participantNicknameHistories := make(map[uuid.UUID][]usecase.ChainedNodeNickname)
	query := `
				WITH RECURSIVE chain_tree AS (
					SELECT 
						c.chain_id,
						c.nickname_id,
						c.parent_chain_id,
						c.chained_at,
						c.active,
						c.chain_id AS root_chain_id,  -- сохраняем корень
						1 AS depth
					FROM chains c
					WHERE c.chain_id = ANY($1::uuid[])
				
					UNION ALL
				
					SELECT 
						c.chain_id,
						c.nickname_id,
						c.parent_chain_id,
						c.chained_at,
						c.active,
						ct.root_chain_id,              -- протягиваем root вниз
						ct.depth + 1
					FROM chains c
					INNER JOIN chain_tree ct ON c.parent_chain_id = ct.chain_id
				)
				SELECT 
					chain_id,
					nickname_id,
					parent_chain_id,
					chained_at,
					active,
					root_chain_id,
					depth
				FROM chain_tree
				ORDER BY root_chain_id, depth;
			 `
	rows, err := q.exec.QueryContext(ctx, query, pq.Array(activeChainIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var node usecase.ChainedNodeNickname
		var depth int
		var rootID uuid.UUID
		err = rows.Scan(
			&node.ChainID,
			&node.NicknameID,
			&node.ParentChainID,
			&node.OwnedAt,
			&node.Active,
			&rootID,
			&depth,
		)
		if err != nil {
			return nil, err
		}
		participantNicknameHistories[rootID] = append(participantNicknameHistories[rootID], node)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	allChainedParticipants := make(usecase.AllChainedParticipants)
	getParticipantQuery := `
								WITH RECURSIVE chain_tree AS (
									-- начальный список цепочек
									SELECT 
										c.chain_id,
										c.nickname_id,
										c.parent_chain_id,
										c.chained_at,
										c.active,
										c.chain_id AS root_chain_id
									FROM chains c
									WHERE c.chain_id = ANY($1::uuid[])
								
									UNION ALL
								
									-- рекурсивное построение вниз
									SELECT 
										c.chain_id,
										c.nickname_id,
										c.parent_chain_id,
										c.chained_at,
										c.active,
										ct.root_chain_id
									FROM chains c
									INNER JOIN chain_tree ct ON c.parent_chain_id = ct.chain_id
								)
								SELECT DISTINCT ON (root_chain_id)
									nickname_id,
									active,
									n.name,
									g.name,
									g.id,
									root_chain_id
								FROM chain_tree
								JOIN aa_nicknames n ON n.id = nickname_id
								JOIN aa_guild_nicknames gn ON gn.nickname_id = n.id
								JOIN aa_guilds g ON g.id = gn.guild_id
								
								WHERE active = TRUE
								ORDER BY root_chain_id DESC;
							`
	rows, err = q.exec.QueryContext(ctx, getParticipantQuery, pq.Array(activeChainIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p usecase.ChainedParticipant
		var active bool
		var rootID uuid.UUID
		err = rows.Scan(
			&p.CurrentNicknameID,
			&active,
			&p.CurrentNickname,
			&p.GuildName,
			&p.GuildID,
			&rootID,
		)
		allChainedParticipants[rootID] = &p
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	for rootID := range participantNicknameHistories {
		if _, ok := allChainedParticipants[rootID]; ok {
			allChainedParticipants[rootID].History = participantNicknameHistories[rootID]
		}
	}
	return allChainedParticipants, nil
}

func (q *GetChainedParticipantsByTenantIDQuery) WithTx(tx *sql.Tx) *GetChainedParticipantsByTenantIDQuery {
	return &GetChainedParticipantsByTenantIDQuery{
		tx,
	}
}
