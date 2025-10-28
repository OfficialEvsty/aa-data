package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type GetTenantRequestQuery struct {
	exec db.ISqlExecutor
}

func NewGetTenantRequestQuery(exec db.ISqlExecutor) *GetTenantRequestQuery {
	return &GetTenantRequestQuery{
		exec: exec,
	}
}

func (q *GetTenantRequestQuery) Handle(ctx context.Context, tenantID uuid.UUID) ([]*domain.Request, error) {
	query := `SELECT id, type, payload, done, created_at, solved_at, rollback_at, is_deleted
			  FROM requests
			  JOIN tenant_requests tr ON tr.request_id = requests.id
			  WHERE tr.tenant_id = $1`
	rows, err := q.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []*domain.Request
	for rows.Next() {
		var request domain.Request
		err = rows.Scan(
			&request.ID,
			&request.Type,
			&request.Payload,
			&request.Done,
			&request.CreatedAt,
			&request.SolvedAt,
			&request.RollbackAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &request)
	}
	return requests, nil
}
