package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetSalaryContextByLunarkIdQuery struct {
	exec db.ISqlExecutor
}

func NewGetSalaryContextByLunarkIdQuery(sql db.ISqlExecutor) *GetSalaryContextByLunarkIdQuery {
	return &GetSalaryContextByLunarkIdQuery{sql}
}

func (q *GetSalaryContextByLunarkIdQuery) Handle(
	ctx context.Context,
	lunarkID uuid.UUID,
) (*usecase.SalaryContext, error) {
	query := `SELECT ls.lunark_id, ls.salary_id, s.fond, s.min_attendance, s.tax, s.created_at, s.status, s.submitted, s.version
 			  FROM lunark_salaries ls
 			  JOIN salaries s ON s.id = ls.salary_id
 			  WHERE ls.lunark_id = $1 AND s.is_deleted = FALSE`
	var result usecase.SalaryContext
	err := q.exec.QueryRowContext(
		ctx,
		query,
		lunarkID,
	).Scan(
		&result.LunarkID,
		&result.SalaryID,
		&result.Fond,
		&result.MinAttendance,
		&result.Tax,
		&result.CreatedAt,
		&result.Status,
		&result.Version,
		&result.SubmittedBy,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
