package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type GetChainIdsByNicknameIdsAtDateIntervalQuery struct {
	exec db.ISqlExecutor
}

func NewGetChainIdsByNicknameIdsAtDateIntervalQuery(sql db.ISqlExecutor) *GetChainIdsByNicknameIdsAtDateIntervalQuery {
	return &GetChainIdsByNicknameIdsAtDateIntervalQuery{sql}
}

// get root chains
func (q *GetChainIdsByNicknameIdsAtDateIntervalQuery) Handle(
	ctx context.Context,
	nicknameIDs []uuid.UUID,
	start time.Time,
	end time.Time,
) (usecase.ChainedNicknames, error) {
	query := `
				WITH RECURSIVE root_chain AS (
					-- стартовые ноды по списку nickname_id
					SELECT 
						c.chain_id,
						c.nickname_id,
						c.parent_chain_id,
						c.chained_at,
						c.active,
						c.nickname_id AS searched_nickname_id
					FROM chains c
					WHERE c.nickname_id = ANY($1::uuid[])
					  AND c.chained_at <= $3
				
					UNION ALL
				
					-- поднимаемся вверх
					SELECT 
						c.chain_id,
						c.nickname_id,
						c.parent_chain_id,
						c.chained_at,
						c.active,
						rc.searched_nickname_id
					FROM chains c
					INNER JOIN root_chain rc ON rc.parent_chain_id = c.chain_id
				),
				intervals AS (
					SELECT 
						c.chain_id,
						c.nickname_id,
						c.parent_chain_id,
						c.chained_at AS start_date,
						COALESCE(child.chained_at, '9999-12-31'::timestamp) AS end_date,
						c.active,
						c.nickname_id AS searched_nickname_id
					FROM chains c
					LEFT JOIN chains child ON child.parent_chain_id = c.chain_id
				)
				SELECT 
					rc.chain_id AS root_chain_id,
					rc.searched_nickname_id,
					i.start_date,
					i.end_date
				FROM root_chain rc
				JOIN intervals i ON i.chain_id = rc.chain_id
				WHERE rc.parent_chain_id IS NULL
				  -- AND i.start_date < i.end_date  -- базовое условие валидного интервала (deprecated)
				  AND i.start_date < $3          -- пересечение с заданным интервалом [$3, $4]
				  AND i.end_date   > $2;	

			`
	rows, err := q.exec.QueryContext(ctx, query, pq.Array(nicknameIDs), start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chainedNicknames = make(usecase.ChainedNicknames)
	for rows.Next() {
		var rootChainID uuid.UUID
		var searchedNicknameID uuid.UUID
		var startDate, endDate time.Time
		err = rows.Scan(&rootChainID, &searchedNicknameID, &startDate, &endDate)
		if err != nil {
			return nil, err
		}
		chainedNicknames[searchedNicknameID] = rootChainID
	}
	return chainedNicknames, nil
}
