package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type TenantLunarkRepository struct {
	exec db.ISqlExecutor
}

func NewTenantLunarkRepository(exec db.ISqlExecutor) *TenantLunarkRepository {
	return &TenantLunarkRepository{exec}
}

func (r *TenantLunarkRepository) Add(ctx context.Context, entry domain.Journal) error {
	query := `INSERT INTO tenant_lunark (tenant_id, lunark_id) VALUES ($1, $2)`
	_, err := r.exec.ExecContext(ctx, query, entry.TenantID, entry.LunarkID)
	if err != nil {
		return err
	}
	return nil
}
func (r *TenantLunarkRepository) Remove(ctx context.Context, lunarkID uuid.UUID) error {
	query := `DELETE FROM tenant_lunark WHERE lunark_id = $1`
	_, err := r.exec.ExecContext(ctx, query, lunarkID)
	if err != nil {
		return err
	}
	return nil
}
func (r *TenantLunarkRepository) All(ctx context.Context, tenantID uuid.UUID) ([]*domain.Journal, error) {
	var result []*domain.Journal
	query := `SELECT tenant_id, lunark_id FROM tenant_lunark`
	rows, err := r.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry domain.Journal
		err = rows.Scan(&entry.TenantID, &entry.LunarkID)
		if err != nil {
			return nil, err
		}
		result = append(result, &entry)
	}
	return result, nil
}

func (r *TenantLunarkRepository) GetByID(ctx context.Context, lunarkID uuid.UUID) (*domain.Journal, error) {
	query := `SELECT tenant_id, lunark_id FROM tenant_lunark WHERE tenant_id = $1`
	row := r.exec.QueryRowContext(ctx, query, lunarkID)
	var entry domain.Journal
	err := row.Scan(&entry.TenantID, &entry.LunarkID)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
func (r *TenantLunarkRepository) WithTx(tx *sql.Tx) junction_repos.ITenantLunarkRepository {
	return &TenantLunarkRepository{tx}
}
