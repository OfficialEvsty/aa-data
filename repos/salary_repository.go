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
	query := `INSERT INTO salaries (id, fond, min_attendance, tax, submitted) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.exec.ExecContext(ctx, query, salary.ID, salary.Fond, salary.MinAttendance, salary.Tax, salary.SubmittedBy)
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
				  tax = $4,
            	  status = $5,
            	  submitted = $6,
              	  version=version+1
			  WHERE id = $1 AND is_deleted = FALSE`
	_, err := r.exec.ExecContext(ctx, query, salary.ID, salary.Fond, salary.MinAttendance, salary.Tax, salary.Status, salary.SubmittedBy)
	return err
}

func (r *SalaryRepository) GetByID(ctx context.Context, salaryID uuid.UUID) (*domain.Salary, error) {
	query := `SELECT id, fond, min_attendance, tax, status, submitted, version FROM salaries WHERE id = $1 AND is_deleted = FALSE`
	row := r.exec.QueryRowContext(ctx, query, salaryID)
	var salary domain.Salary
	err := row.Scan(
		&salary.ID,
		&salary.Fond,
		&salary.MinAttendance,
		&salary.Tax,
		&salary.Status,
		&salary.SubmittedBy,
		&salary.Version,
	)
	return &salary, err
}
func (r *SalaryRepository) WithTx(tx *sql.Tx) repos.ISalaryRepository {
	return &SalaryRepository{tx}
}
