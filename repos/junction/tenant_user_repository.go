package junction_repos

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type TenantUserRepository struct {
	exec db.ISqlExecutor
}

func NewTenantUserRepository(exec db.ISqlExecutor) *TenantUserRepository {
	return &TenantUserRepository{exec}
}

func (r *TenantUserRepository) Add(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) error {
	query := `INSERT INTO tenant_users (tenant_id, user_id) VALUES ($1, $2) ON CONFLICT (tenant_id, user_id) DO NOTHING`
	_, err := r.exec.ExecContext(ctx, query, tenantID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (r *TenantUserRepository) Remove(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM tenant_users WHERE tenant_id = $1 AND user_id = $2`
	_, err := r.exec.ExecContext(ctx, query, tenantID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TenantUserRepository) GetTenant(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	query := `SELECT tenant_id FROM tenant_users WHERE user_id = $1`
	row := r.exec.QueryRowContext(ctx, query, userID)
	var tenantID uuid.UUID
	err := row.Scan(&tenantID)
	if err != nil {
		return uuid.Nil, err
	}
	return tenantID, nil
}
func (r *TenantUserRepository) GetUserIDs(ctx context.Context, tenantID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM tenant_users WHERE tenant_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}
func (r *TenantUserRepository) CheckUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `SELECT tenant_id FROM tenant_users WHERE tenant_id = $1 AND user_id = $2`
	row := r.exec.QueryRowContext(ctx, query, tenantID, userID)
	err := row.Scan(&tenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *TenantUserRepository) WithTx(tx *sql.Tx) junction_repos.ITenantUserRepository {
	return &TenantUserRepository{tx}
}
