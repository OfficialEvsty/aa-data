package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/errors"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TenantChainRepository struct {
	exec db.ISqlExecutor
}

func NewTenantChainRepository(exec db.ISqlExecutor) *TenantChainRepository {
	return &TenantChainRepository{exec}
}

func (r *TenantChainRepository) Add(ctx context.Context, tc domain.TenantChain) error {
	query := `INSERT INTO tenant_chains (tenant_id, chain_id) VALUES ($1, $2) ON CONFLICT (tenant_id, chain_id) DO NOTHING`
	_, err := r.exec.ExecContext(ctx, query, tc.TenantID, tc.ChainID)
	return err
}

func (r *TenantChainRepository) GetAllByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantChain, error) {
	var result []*domain.TenantChain
	query := `SELECT tenant_id, chain_id FROM tenant_chains WHERE tenant_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tc := &domain.TenantChain{}
		err = rows.Scan(&tc.TenantID, &tc.ChainID)
		if err != nil {
			return nil, err
		}
		result = append(result, tc)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *TenantChainRepository) Remove(ctx context.Context, chainID uuid.UUID) error {
	query := `DELETE FROM tenant_chains WHERE chain_id = $1`
	_, err := r.exec.ExecContext(ctx, query, chainID)
	return err
}

func (r *TenantChainRepository) CheckTenantAttachment(ctx context.Context, tenantID uuid.UUID, chainIDs []uuid.UUID) error {
	query := `SELECT tenant_id, chain_id FROM tenant_chains WHERE tenant_id = $1 AND chain_id = ANY($2)`
	rows, err := r.exec.QueryContext(ctx, query, tenantID, pq.Array(chainIDs))
	if err != nil {
		return err
	}
	defer rows.Close()
	var count int = 0
	for rows.Next() {
		tc := &domain.TenantChain{}
		err = rows.Scan(&tc.TenantID, &tc.ChainID)
		if err != nil {
			return err
		}
		count++
	}
	if count != len(chainIDs) {
		return errors.ErrorNotAttachedToSpecifiedTenant
	}
	return nil
}

func (r *TenantChainRepository) WithTx(tx *sql.Tx) junction_repos.ITenantChainRepository {
	return &TenantChainRepository{tx}
}
