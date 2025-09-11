package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos2 "github.com/OfficialEvsty/aa-data/repos"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/OfficialEvsty/aa-shared/types"
	"github.com/google/uuid"
)

// NicknameHistory stack where root nickname is tail element and active is head
type NicknameHistory = types.Stack[*domain.NicknameChain]

type GetNicknamesHistoryByRelativeChainIdQuery struct {
	exec      db.ISqlExecutor
	chainRepo repos.IChainRepository
}

func NewGetNicknamesHistoryByRelativeChainIdQuery(exec db.ISqlExecutor) *GetNicknamesHistoryByRelativeChainIdQuery {
	return &GetNicknamesHistoryByRelativeChainIdQuery{
		exec:      exec,
		chainRepo: repos2.NewChainRepository(exec),
	}
}

func (q *GetNicknamesHistoryByRelativeChainIdQuery) Handle(ctx context.Context, chainID uuid.UUID) (*NicknameHistory, error) {
	// Part 1: Receiving root chain id
	rootChain, err := q.chainRepo.GetRootChain(ctx, chainID)
	if err != nil {
		return nil, err
	}
	// Part 2: Building chain's history by root chain id
	var history NicknameHistory
	query := `
				WITH RECURSIVE chain_tree AS (
					SELECT chain_id, nickname_id, parent_chain_id, chained_at, active, 1 AS depth
					FROM chains
					WHERE chain_id = $1
				
					UNION ALL
					
					SELECT c.chain_id, c.nickname_id, c.parent_chain_id, c.chained_at, c.active, ct.depth + 1
					FROM chains c 
					INNER JOIN chain_tree ct ON c.parent_chain_id = ct.chain_id
				)
				SELECT chain_id, nickname_id, parent_chain_id, chained_at, active
				FROM chain_tree
				ORDER BY depth;
			 `
	rows, err := q.exec.QueryContext(ctx, query, rootChain.ChainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var chain domain.NicknameChain
		err = rows.Scan(
			&chain.ChainID,
			&chain.NicknameID,
			&chain.ParentChainID,
			&chain.ChainedAt,
			&chain.Active,
		)
		history.Push(&chain)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &history, nil
}

func (q *GetNicknamesHistoryByRelativeChainIdQuery) WithTx(tx *sql.Tx) *GetNicknamesHistoryByRelativeChainIdQuery {
	return &GetNicknamesHistoryByRelativeChainIdQuery{tx, q.chainRepo.WithTx(tx)}
}
