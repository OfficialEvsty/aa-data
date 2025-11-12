package junction_repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
)

type ArchivedChainRepository struct {
	db db.ISqlExecutor
}

func NewArchivedChainRepository(db db.ISqlExecutor) *ArchivedChainRepository {
	return &ArchivedChainRepository{db: db}
}

func (r *ArchivedChainRepository) Add(ctx context.Context, aChain domain.ArchivedChain) error {
	query := `INSERT INTO archived_chains (tenant_id, chain_id) VALUES ($1, $2) ON CONFLICT (chain_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, aChain.TenantID, aChain.ChainID)
	return err
}
func (r *ArchivedChainRepository) AddMany(ctx context.Context, aChains []domain.ArchivedChain) error {
	argCount := 2
	valueStrings := make([]string, 0, len(aChains))
	valueArgs := make([]interface{}, 0, len(aChains)*argCount)

	for i, p := range aChains {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, p.TenantID, p.ChainID)
	}

	stmt := fmt.Sprintf("INSERT INTO archived_chains (tenant_id, chain_id) VALUES %s ON CONFLICT (chain_id) DO NOTHING", strings.Join(valueStrings, ","))
	_, err := r.db.ExecContext(ctx, stmt, valueArgs...)
	return err
}
func (r *ArchivedChainRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.ArchivedChain, error) {
	var res []*domain.ArchivedChain
	query := `SELECT tenant_id, chain_id, archived_at FROM archived_chains WHERE tenant_id = $1`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var aChain domain.ArchivedChain
		err = rows.Scan(
			&aChain.TenantID,
			&aChain.ChainID,
			&aChain.ArchivedAt,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, &aChain)
	}
	return res, nil
}
func (r *ArchivedChainRepository) IsArchived(ctx context.Context, chainID uuid.UUID) (bool, error) {
	query := `SELECT chain_id FROM archived_chains WHERE chain_id = $1`
	row := r.db.QueryRowContext(ctx, query, chainID)
	err := row.Scan(&chainID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (r *ArchivedChainRepository) Remove(ctx context.Context, chainID uuid.UUID) error {
	query := `DELETE FROM archived_chains WHERE chain_id = $1`
	_, err := r.db.ExecContext(ctx, query, chainID)
	return err
}
func (r *ArchivedChainRepository) RemoveMany(ctx context.Context, chainIDs []uuid.UUID) error {
	query := `DELETE FROM archived_chains WHERE chain_id = ANY($1)`
	_, err := r.db.ExecContext(ctx, query, pq.Array(chainIDs))
	return err
}
