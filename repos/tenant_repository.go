package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

type TenantRepository struct {
	exec db.ISqlExecutor
}

func NewTenantRepository(exec db.ISqlExecutor) *TenantRepository {
	return &TenantRepository{exec}
}

func (r *TenantRepository) Add(ctx context.Context, tenant domain.Tenant) (*domain.Tenant, error) {
	query := `INSERT INTO tenants (id, name, owner_id) VALUES ($1, $2, $3) RETURNING id, name, created_at, owner_id;`
	row := r.exec.QueryRowContext(ctx, query, tenant.ID, tenant.Name, tenant.OwnerID)
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, &tenant.OwnerID)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}
func (r *TenantRepository) Remove(ctx context.Context, tenantID uuid.UUID) error {
	query := `DELETE FROM tenants WHERE id = $1;`
	_, err := r.exec.ExecContext(ctx, query, tenantID)
	if err != nil {
		return err
	}
	return nil
}
func (r *TenantRepository) Update(ctx context.Context, tenant domain.Tenant) (*domain.Tenant, error) {
	query := `UPDATE tenants SET name = $2, owner_id = $3 WHERE id = $1 RETURNING id, name, created_at, owner_id;`
	err := r.exec.QueryRowContext(ctx, query, tenant.ID, tenant.Name, tenant.OwnerID).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.CreatedAt,
		&tenant.OwnerID,
	)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}
func (r *TenantRepository) GetByID(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	query := `SELECT id, name, created_at, owner_id FROM tenants WHERE id = $1;`
	row := r.exec.QueryRowContext(ctx, query, tenantID)
	err := row.Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.CreatedAt,
		&tenant.OwnerID,
	)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}
func (r *TenantRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	query := `SELECT id, name, created_at, owner_id FROM tenants WHERE owner_id = $1;`
	row := r.exec.QueryRowContext(ctx, query, ownerID)
	err := row.Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.CreatedAt,
		&tenant.OwnerID,
	)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) GetOwnerID(ctx context.Context, tenantID uuid.UUID) (uuid.UUID, error) {
	var ownerID uuid.UUID
	query := `SELECT owner_id FROM tenants WHERE id = $1;`
	row := r.exec.QueryRowContext(ctx, query, tenantID)
	err := row.Scan(&ownerID)
	if err != nil {
		return uuid.Nil, err
	}
	return ownerID, nil
}

func (r *TenantRepository) All(ctx context.Context) ([]*domain.Tenant, error) {
	var tenants []*domain.Tenant = make([]*domain.Tenant, 0)
	query := `SELECT id, name, created_at, owner_id FROM tenants;`
	rows, err := r.exec.QueryContext(ctx, query)
	if err != nil {
		return tenants, err
	}
	defer rows.Close()
	for rows.Next() {
		var tenant domain.Tenant
		err = rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.CreatedAt,
			&tenant.OwnerID,
		)
		tenants = append(tenants, &tenant)
	}
	err = rows.Err()
	if err != nil {
		return tenants, err
	}
	return tenants, nil
}

func (r *TenantRepository) WithTx(tx *sql.Tx) repos.ITenantRepository {
	return &TenantRepository{tx}
}
