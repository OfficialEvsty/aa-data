package junction_repos

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type LunarSalaryRepository struct {
	exec db.ISqlExecutor
}

func NewLunarkSalaryRepository(exec db.ISqlExecutor) *LunarSalaryRepository {
	return &LunarSalaryRepository{exec: exec}
}

func (r *LunarSalaryRepository) Add(ctx context.Context, lunarkSalary domain.LunarkSalary) error {
	query := `INSERT INTO lunark_salaries (lunark_id, salary_id) VALUES ($1, $2)`
	_, err := r.exec.ExecContext(ctx, query, lunarkSalary.LunarkID, lunarkSalary.SalaryID)
	return err
}
func (r *LunarSalaryRepository) Remove(ctx context.Context, salaryID uuid.UUID) error {
	query := `DELETE FROM lunark_salaries WHERE salary_id = $1`
	_, err := r.exec.ExecContext(ctx, query, salaryID)
	return err
}
func (r *LunarSalaryRepository) GetSalaries(ctx context.Context, lunarkID uuid.UUID) ([]*domain.LunarkSalary, error) {
	var result []*domain.LunarkSalary
	query := `SELECT lunark_id, salary_id FROM lunark_salaries WHERE lunark_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, lunarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var lunarkSalary domain.LunarkSalary
		err = rows.Scan(&lunarkID, &lunarkSalary.SalaryID)
		if err != nil {
			return nil, err
		}
		result = append(result, &lunarkSalary)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}
