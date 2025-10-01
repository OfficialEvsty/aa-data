package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
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

type LunarkPayout struct {
	LunarkID  uuid.UUID                  `json:"lunark_id"`
	Name      string                     `json:"name"`
	StartDate time.Time                  `json:"start_date"`
	EndDate   *time.Time                 `json:"end_date"`
	Opened    bool                       `json:"opened"`
	PayoutID  *uuid.UUID                 `json:"payout_id"`
	Status    *serializable.SalaryStatus `json:"status"`
}

type GetOpenedLunarkByTenant struct {
	exec db.ISqlExecutor
}

func NewGetOpenedLunarkByTenant(exec db.ISqlExecutor) GetOpenedLunarkByTenant {
	return GetOpenedLunarkByTenant{exec: exec}
}

func (q *GetOpenedLunarkByTenant) GetLunarkPayoutList(ctx context.Context, tenantID uuid.UUID) ([]*LunarkPayout, error) {
	var result []*LunarkPayout
	query := `SELECT l.id, l.name, l.start_date, l.end_date, tl.opened, s.id, s.status
              FROM tenant_lunark tl
              JOIN lunark l ON l.id = tl.lunark_id
              LEFT JOIN lunark_salaries ls ON ls.lunark_id = l.id
              LEFT JOIN salaries s ON s.id = ls.salary_id AND s.is_deleted = FALSE
              WHERE tl.tenant_id
              ORDER BY l.start_date DESC`
	rows, err := q.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var l LunarkPayout
		err = rows.Scan(
			&l.LunarkID,
			&l.Name,
			&l.StartDate,
			&l.EndDate,
			&l.Opened,
			&l.PayoutID,
			&l.Status,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &l)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *GetOpenedLunarkByTenant) Handle(ctx context.Context, tenantID uuid.UUID) (*OpenedLunark, error) {
	var opened OpenedLunark
	query := `SELECT tl.tenant_id, tl.lunark_id, l.name, l.start_date, l.end_date, tl.opened
			  FROM tenant_lunark AS tl
			  JOIN lunark AS l ON l.id = tl.lunark_id
			  WHERE tl.tenant_id = $1 AND tl.opened = TRUE`

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
