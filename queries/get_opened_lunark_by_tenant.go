package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/google/uuid"
	"time"
)

type OpenedLunark struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	LunarkID  uuid.UUID  `json:"lunark_id"`
	Name      string     `json:"name"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Opened    bool       `json:"opened"`
}

type GetOpenedLunarkByTenant struct {
	exec db.ISqlExecutor
}

func NewGetOpenedLunarkByTenant(exec db.ISqlExecutor) GetOpenedLunarkByTenant {
	return GetOpenedLunarkByTenant{exec: exec}
}

func (q *GetOpenedLunarkByTenant) Handle(ctx context.Context, tenantID uuid.UUID) (*OpenedLunark, error) {
	var opened OpenedLunark
	query := `SELECT tl.tenant_id, tl.lunark_id, l.name, l.start_date, l.end_date, l.opened
			  FROM tenant_lunark AS tl
			  JOIN lunark AS l ON l.id = tl.lunark_id
			  WHERE tl.tenant_id = $1 AND l.opened = TRUE`

	err := q.exec.QueryRowContext(ctx, query, tenantID).Scan(
		&opened.TenantID,
		&opened.LunarkID,
		&opened.Name,
		&opened.StartDate,
		&opened.EndDate,
		&opened.Opened,
	)
	if err != nil {
		return nil, err
	}
	return &opened, nil
}

func (q *GetOpenedLunarkByTenant) WithTx(tx *sql.Tx) *GetOpenedLunarkByTenant {
	return &GetOpenedLunarkByTenant{
		tx,
	}
}
