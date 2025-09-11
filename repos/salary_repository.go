package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

type SalaryRepository struct {
	exec db.ISqlExecutor
}

func NewSalaryRepository(exec db.ISqlExecutor) *SalaryRepository {
	return &SalaryRepository{exec}
}

func (r *SalaryRepository) Add(ctx context.Context, salary domain.Salary) error {
	query := `INSERT INTO salaries (id, fond, min_attendance, tax) VALUES ($1, $2, $3, $4)`
	_, err := r.exec.ExecContext(ctx, query, salary.ID, salary.Fond, salary.MinAttendance, salary.Tax)
	return err
}
func (r *SalaryRepository) Remove(ctx context.Context, salaryID uuid.UUID) error {
	query := `DELETE FROM salaries WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, salaryID)
	return err
}

func (r *SalaryRepository) Update(ctx context.Context, salary domain.Salary) error {
	query := `UPDATE salaries
              SET 
                  fond = $2,
				  min_attendance = $3, 
				  tax = $4
			  WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, salary.ID, salary.Fond, salary.MinAttendance, salary.Tax)
	return err
}

func (r *SalaryRepository) GetByID(ctx context.Context, salaryID uuid.UUID) (*domain.Salary, error) {
	query := `SELECT id, fond, min_attendance, tax FROM salaries WHERE id = $1`
	row := r.exec.QueryRowContext(ctx, query, salaryID)
	var salary domain.Salary
	err := row.Scan(
		&salary.ID,
		&salary.Fond,
		&salary.MinAttendance,
		&salary.Tax,
	)
	return &salary, err
}
func (r *SalaryRepository) WithTx(tx *sql.Tx) repos.ISalaryRepository {
	return &SalaryRepository{tx}
}
