package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type TenantPublishRepository struct {
	exec db.ISqlExecutor
}

func NewTenantPublishRepository(exec db.ISqlExecutor) *TenantPublishRepository {
	return &TenantPublishRepository{exec}
}

func (r *TenantPublishRepository) Add(ctx context.Context, publish domain.TenantPublish) (*domain.TenantPublish, error) {
	query := `INSERT INTO tenant_publishes (tenant_id, publish_id, user_id)
			  VALUES ($1, $2, $3) RETURNING tenant_id, publish_id, user_id, published_at;`
	err := r.exec.QueryRowContext(ctx, query, publish.TenantID, publish.PublishID, publish.UserID).Scan(&publish.TenantID, &publish.PublishID, &publish.UserID, &publish.PublishedAt)
	if err != nil {
		return nil, err
	}
	return &publish, nil
}

func (r *TenantPublishRepository) Remove(ctx context.Context, publishID uuid.UUID) error {
	query := `DELETE FROM tenant_publishes WHERE publish_id = $1;`
	_, err := r.exec.ExecContext(ctx, query, publishID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TenantPublishRepository) GetByID(ctx context.Context, publishID uuid.UUID) (*domain.TenantPublish, error) {
	var TenantPublish domain.TenantPublish
	query := `SELECT tenant_id, publish_id, user_id, published_at FROM tenant_publishes WHERE publish_id = $1;`
	err := r.exec.QueryRowContext(ctx, query, publishID).Scan(&TenantPublish.TenantID, &TenantPublish.PublishID, &TenantPublish.UserID, &TenantPublish.PublishedAt)
	if err != nil {
		return nil, err
	}
	return &TenantPublish, nil
}

func (r *TenantPublishRepository) All(ctx context.Context, TenantID uuid.UUID) ([]*domain.TenantPublish, error) {
	query := `SELECT tenant_id, publish_id, user_id, published_at FROM tenant_publishes WHERE tenant_id = $1;`
	rows, err := r.exec.QueryContext(ctx, query, TenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var TenantPublishes []*domain.TenantPublish
	for rows.Next() {
		var TenantPublish domain.TenantPublish
		err = rows.Scan(&TenantPublish.TenantID, &TenantPublish.PublishID, &TenantPublish.UserID, &TenantPublish.PublishedAt)
		if err != nil {
			return nil, err
		}
		TenantPublishes = append(TenantPublishes, &TenantPublish)
	}
	return TenantPublishes, nil
}
func (r *TenantPublishRepository) RemoveAll(ctx context.Context, TenantID uuid.UUID) error {
	return nil
}
func (r *TenantPublishRepository) WithTx(tx *sql.Tx) junction_repos.ITenantPublishRepository {
	return &TenantPublishRepository{tx}
}
