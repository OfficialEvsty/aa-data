package junction_repos

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
	"strings"
)

type PaymentRepository struct {
	exec db.ISqlExecutor
}

func NewPaymentRepository(exec db.ISqlExecutor) *PaymentRepository {
	return &PaymentRepository{exec}
}

func (r *PaymentRepository) AddMany(ctx context.Context, payments []domain.Payment) error {
	valueStrings := make([]string, 0, len(payments))
	valueArgs := make([]interface{}, 0, len(payments)*4)

	for i, p := range payments {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, p.SalaryID, p.ChainID, p.Value, p.Reason)
	}

	stmt := fmt.Sprintf("INSERT INTO payments (salary_id, root_chain_id, salary, reason) VALUES %s ON CONFLICT (salary_id, root_chain_id) DO UPDATE SET reason=EXCLUDED.reason, salary=EXCLUDED.salary", strings.Join(valueStrings, ","))
	_, err := r.exec.ExecContext(ctx, stmt, valueArgs...)
	return err
}
func (r *PaymentRepository) Clear(ctx context.Context, salaryID uuid.UUID) error {
	stmt := fmt.Sprintf("DELETE FROM payments WHERE salary_id = $1", salaryID)
	_, err := r.exec.ExecContext(ctx, stmt)
	return err
}
func (r *PaymentRepository) GetAllBySalaryID(ctx context.Context, salaryID uuid.UUID) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	query := `SELECT salary_id, root_chain_id, salary, reason FROM payments WHERE salary_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, salaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var payment domain.Payment
		rows.Scan(
			&payment.SalaryID,
			&payment.ChainID,
			&payment.Value,
			&payment.Reason,
		)
		payments = append(payments, &payment)
	}
	return payments, nil
}
func (r *PaymentRepository) GetAllByChainID(ctx context.Context, chainID uuid.UUID) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	query := `SELECT salary_id, root_chain_id, salary, reason FROM payments WHERE salary_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var payment domain.Payment
		err = rows.Scan(
			&payment.SalaryID,
			&payment.ChainID,
			&payment.Value,
			&payment.Reason,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}
	return payments, nil
}
func (r *PaymentRepository) WithTx(tx *sql.Tx) junction_repos.IPaymentRepository {
	return &PaymentRepository{tx}
}
