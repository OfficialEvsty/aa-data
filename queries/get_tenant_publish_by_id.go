package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
)

// TenantPublishDTO is result of GetTenantPublishByIDQuery
type TenantPublishDTO struct {
	PublishID  uuid.UUID                 `json:"publish_id"`
	UserID     uuid.UUID                 `json:"user_id"`
	TenantID   uuid.UUID                 `json:"tenant_id"`
	TenantName string                    `json:"tenant_name"`
	S3         serializable.S3Screenshot `json:"s3"`
}

type GetTenantPublishByIDQuery struct {
	exec db.ISqlExecutor
}

func NewGetTenantPublishByIDQuery(exec db.ISqlExecutor) *GetTenantPublishByIDQuery {
	return &GetTenantPublishByIDQuery{
		exec: exec,
	}
}

func (q GetTenantPublishByIDQuery) Handle(ctx context.Context, publishID uuid.UUID) (*TenantPublishDTO, error) {
	var dto TenantPublishDTO
	query := `SELECT p.id, t.id, t.name, tp.user_id, p.s3 FROM tenant_publishes as tp
			  JOIN publishes as p ON p.id = tp.publish_id
			  JOIN tenants as t ON t.id = tp.tenant_id
			  WHERE tp.publish_id = $1`
	err := q.exec.QueryRowContext(ctx, query, publishID).Scan(&dto.PublishID, &dto.TenantID, &dto.TenantName, &dto.UserID, &dto.S3)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}
