package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/google/uuid"
	"time"
)

type TenantUser struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Username string    `json:"username"`
	JoinedAt time.Time `json:"joined_at"`
	Email    string    `json:"email"`
	LastSeen time.Time `json:"last_seen"`
}

type GetTenantUserByTenantID struct {
	exec db.ISqlExecutor
}

func NewGetTenantUserByTenantID(exec db.ISqlExecutor) *GetTenantUserByTenantID {
	return &GetTenantUserByTenantID{
		exec: exec,
	}
}

func (q *GetTenantUserByTenantID) Handle(ctx context.Context, tenantID uuid.UUID) ([]*TenantUser, error) {
	var res []*TenantUser
	query := `SELECT u.id, tu.tenant_id, u.username, u.email, u.created_at, u.last_seen FROM tenant_users tu
			  JOIN users as u ON u.id = tu.user_id
			  WHERE tu.tenant_id = $1`
	rows, err := q.exec.QueryContext(
		ctx,
		query,
		tenantID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var dto TenantUser
		err = rows.Scan(
			&dto.UserID,
			&dto.TenantID,
			&dto.Username,
			&dto.Email,
			&dto.JoinedAt,
			&dto.LastSeen,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, &dto)
	}
	return res, nil
}
