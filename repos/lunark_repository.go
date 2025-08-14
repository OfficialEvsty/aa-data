package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
	"time"
)

type LunarkRepository struct {
	exec db.ISqlExecutor
}

func NewLunarkRepository(exec db.ISqlExecutor) *LunarkRepository {
	return &LunarkRepository{exec}
}

func (r *LunarkRepository) Add(ctx context.Context, lunark domain.Lunark) error {
	query := `INSERT INTO lunark (id, name, start_date) VALUES ($1, $2, $3)`
	_, err := r.exec.ExecContext(ctx, query, lunark.ID, lunark.Name, lunark.StartDate)
	if err != nil {
		return err
	}
	return nil
}
func (r *LunarkRepository) Update(ctx context.Context, lunark domain.Lunark) error {
	query := `UPDATE lunark SET name = $2, start_date = $3, end_date = $4 WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, lunark.ID, lunark.Name, lunark.StartDate, lunark.EndDate)
	if err != nil {
		return err
	}
	return nil
}

func (r *LunarkRepository) UpdateEndDate(ctx context.Context, id uuid.UUID, end time.Time) error {
	query := `UPDATE lunark SET end_date = $2 WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, id, end)
	if err != nil {
		return err
	}
	return nil
}

func (r *LunarkRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Lunark, error) {
	var lunark domain.Lunark
	query := `SELECT id, name, start_date, end_date FROM lunark WHERE id = $1`
	row := r.exec.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&lunark.ID,
		&lunark.Name,
		&lunark.StartDate,
		&lunark.EndDate,
	)
	if err != nil {
		return nil, err
	}
	return &lunark, nil
}
func (r *LunarkRepository) WithTx(tx *sql.Tx) repos.ILunarkRepository {
	return &LunarkRepository{tx}
}
