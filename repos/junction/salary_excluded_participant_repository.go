package junction_repos

import (
	"context"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
	"strings"
)

type SalaryExcludedParticipantRepository struct {
	exec db.ISqlExecutor
}

func NewSalaryExcludedParticipantRepository(exec db.ISqlExecutor) *SalaryExcludedParticipantRepository {
	return &SalaryExcludedParticipantRepository{exec: exec}
}

func (r *SalaryExcludedParticipantRepository) AddMany(ctx context.Context, excludes []domain.ExcludedParticipant) error {
	valueStrings := make([]string, 0, len(excludes))
	valueArgs := make([]interface{}, 0, len(excludes)*3)

	for i, exclude := range excludes {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, exclude.SalaryID, exclude.ChainID, exclude.Reason)
	}

	stmt := fmt.Sprintf("INSERT INTO excluded_participants_salary (salary_id, root_chain_id, reason) VALUES %s ON CONFLICT (salary_id, root_chain_id) DO UPDATE SET reason=EXCLUDED.reason", strings.Join(valueStrings, ","))
	_, err := r.exec.ExecContext(ctx, stmt, valueArgs...)
	return err
}
func (r *SalaryExcludedParticipantRepository) Clear(ctx context.Context, salaryID uuid.UUID) error {
	query := `DELETE FROM excluded_participants_salary WHERE salary_id = $1`
	_, err := r.exec.ExecContext(ctx, query, salaryID)
	return err
}
