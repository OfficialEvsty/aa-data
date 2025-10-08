package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/google/uuid"
	"time"
)

type TenantChainDTO struct {
	ChainID    uuid.UUID `json:"chain_id"`
	ChainedAt  time.Time `json:"chained_at"`
	NicknameID uuid.UUID `json:"nickname_id"`
	Active     bool      `json:"active"`
}

type GetActiveChainsByTenantQuery struct {
	exec db.ISqlExecutor
}

func NewGetActiveChainsByTenantQuery(exec db.ISqlExecutor) *GetActiveChainsByTenantQuery {
	return &GetActiveChainsByTenantQuery{exec: exec}
}

func (q *GetActiveChainsByTenantQuery) Handle(ctx context.Context, tenantID uuid.UUID) ([]*TenantChainDTO, error) {
	var result []*TenantChainDTO = make([]*TenantChainDTO, 0)
	query := `SELECT c.chain_id, c.nickname_id, c.active, c.chained_at
              FROM chains c
              JOIN tenant_chains tc ON c.chain_id = tc.chain_id
              WHERE tc.tenant_id = $1 AND c.active = TRUE`
	rows, err := q.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var t TenantChainDTO
		err = rows.Scan(&t.ChainID, &t.NicknameID, &t.Active, &t.ChainedAt)
		if err != nil {
			return result, err
		}
		result = append(result, &t)
	}
	err = rows.Err()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (q *GetActiveChainsByTenantQuery) GetRootChainIDs(ctx context.Context, tenantID uuid.UUID) ([]*TenantChainDTO, error) {
	var result []*TenantChainDTO = make([]*TenantChainDTO, 0)
	query := `SELECT c.chain_id, c.nickname_id, c.active, c.chained_at
              FROM chains c
              JOIN tenant_chains tc ON c.chain_id = tc.chain_id
              WHERE tc.tenant_id = $1 AND c.parent_chain_id IS NULL`
	rows, err := q.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var t TenantChainDTO
		err = rows.Scan(&t.ChainID, &t.NicknameID, &t.Active, &t.ChainedAt)
		if err != nil {
			return result, err
		}
		result = append(result, &t)
	}
	err = rows.Err()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (q *GetActiveChainsByTenantQuery) WithTx(tx *sql.Tx) *GetActiveChainsByTenantQuery {
	return &GetActiveChainsByTenantQuery{tx}
}
