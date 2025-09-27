package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type ExcludedParticipantRepository struct {
	exec db.ISqlExecutor
}

func NewExcludedParticipantRepository(exec db.ISqlExecutor) *ExcludedParticipantRepository {
	return &ExcludedParticipantRepository{exec}
}

func (r *ExcludedParticipantRepository) Add(ctx context.Context, exP usecase.ExcludedParticipant) error {
	query := `INSERT INTO excluded_participants_salary
              (root_chain_id, salary_id, reason) VALUES ($1, $2, $3)`
	_, err := r.exec.ExecContext(ctx, query, exP.ChainID, exP.SalaryID, exP.Reason)
	return err
}
func (r *ExcludedParticipantRepository) Remove(ctx context.Context, exPID uuid.UUID, salaryID uuid.UUID) error {
	query := `DELETE FROM excluded_participants_salary WHERE salary_id = $1 AND root_chain_id = $2`
	_, err := r.exec.ExecContext(ctx, query, salaryID, exPID)
	return err
}
func (r *ExcludedParticipantRepository) All(ctx context.Context, salaryID uuid.UUID) ([]*usecase.ExcludedParticipant, error) {
	query := `SELECT root_chain_id, salary_id, reason 
			  FROM excluded_participants_salary
			  WHERE salary_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, salaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exPs := make([]*usecase.ExcludedParticipant, 0)
	for rows.Next() {
		var exP usecase.ExcludedParticipant
		err = rows.Scan(&exP.ChainID, &exP.SalaryID, &exP.Reason)
		if err != nil {
			return nil, err
		}
		exPs = append(exPs, &exP)
	}
	return exPs, nil
}
func (r *ExcludedParticipantRepository) WithTx(tx *sql.Tx) junction_repos.IExcludedParticipantRepository {
	return &ExcludedParticipantRepository{tx}
}
