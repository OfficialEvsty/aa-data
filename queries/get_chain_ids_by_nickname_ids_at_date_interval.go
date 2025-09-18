package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
	"time"
)

type GetChainIdsByNicknameIdsAtDateIntervalQuery struct {
	exec db.ISqlExecutor
}

func NewGetChainIdsByNicknameIdsAtDateIntervalQuery(sql db.ISqlExecutor) *GetChainIdsByNicknameIdsAtDateIntervalQuery {
	return &GetChainIdsByNicknameIdsAtDateIntervalQuery{sql}
}

func (q *GetChainIdsByNicknameIdsAtDateIntervalQuery) Handle(
	ctx context.Context,
	nicknameIDs []uuid.UUID,
	start time.Time,
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
					  AND c.chained_at <= $2
				
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
				)
				SELECT 
					chain_id AS root_chain_id,
					searched_nickname_id
				FROM root_chain
				WHERE parent_chain_id IS NULL;
			`
	rows, err := q.exec.QueryContext(ctx, query, nicknameIDs, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chainedNicknames = make(usecase.ChainedNicknames)
	for rows.Next() {
		var rootChainID uuid.UUID
		var searchedNicknameID uuid.UUID
		err = rows.Scan(&rootChainID, &searchedNicknameID)
		if err != nil {
			return nil, err
		}
		chainedNicknames[searchedNicknameID] = rootChainID
	}
	return chainedNicknames, nil
}
