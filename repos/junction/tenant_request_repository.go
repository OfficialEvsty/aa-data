package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type TenantRequestRepository struct {
	db db.ISqlExecutor
}

func NewTenantRequestRepository(db db.ISqlExecutor) *TenantRequestRepository {
	return &TenantRequestRepository{
		db: db,
	}
}

func (r *TenantRequestRepository) Add(ctx context.Context, tenantRequest domain.TenantRequest) error {
	query := `INSERT INTO tenant_requests (tenant_id, request_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, tenantRequest.TenantID, tenantRequest.RequestID)
	return err
}

func (r *TenantRequestRepository) Remove(ctx context.Context, requestID uuid.UUID) error {
	query := `DELETE FROM tenant_requests WHERE request_id = $1`
	_, err := r.db.ExecContext(ctx, query, requestID)
	return err
}
func (r *TenantRequestRepository) GetAllByTenantID(ctx context.Context, tenantID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT request_id FROM tenant_requests WHERE tenant_id = $1`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	requestIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var requestID uuid.UUID
		if err := rows.Scan(&requestID); err != nil {
			return nil, err
		}
		requestIDs = append(requestIDs, requestID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return requestIDs, nil
}

func (r *TenantRequestRepository) WithTx(tx *sql.Tx) junction_repos.ITenantRequestRepository {
	return &TenantRequestRepository{
		db: tx,
	}
}
