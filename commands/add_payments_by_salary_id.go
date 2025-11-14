package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	errors2 "github.com/OfficialEvsty/aa-data/errors"
	"github.com/OfficialEvsty/aa-data/repos"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
	"math"
)

type AddPaymentsBySalaryIDCommand struct {
	PaymentsByChainID []domain.Payment
	SalaryID          uuid.UUID
	CurrentVersion    int
}

type PaymentManager struct {
	txManager   *db.TxManager
	paymentRepo junction_repos.IPaymentRepository
	salaryRepo  repos2.ISalaryRepository
}

func NewPaymentManager(sql *sql.DB) *PaymentManager {
	return &PaymentManager{
		txManager:   db.NewTxManager(sql),
		paymentRepo: junction_repos2.NewPaymentRepository(sql),
		salaryRepo:  repos.NewSalaryRepository(sql),
	}
}

func (m *PaymentManager) Handle(
	ctx context.Context,
	command *AddPaymentsBySalaryIDCommand,
) error {
	return m.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		salary, err := m.salaryRepo.WithTx(tx).GetByID(ctx, command.SalaryID)
		if err != nil {
			return err
		}
		if salary.Version != command.CurrentVersion {
			return errors2.ErrorSalaryVersionMismatch
		}

		err = m.paymentRepo.WithTx(tx).AddMany(ctx, command.PaymentsByChainID)
		if err != nil {
			return err
		}
		salary.Status = serializable.SalaryInProcessing
		err = m.salaryRepo.WithTx(tx).Update(ctx, *salary)
		if err != nil {
			return err
		}
		query := `SELECT SUM(p.salary), s.fond, s.tax FROM payments p
				  JOIN salaries s ON s.id = p.salary_id
				  WHERE salary_id = $1
				  GROUP BY s.fond, s.tax;`
		var salarySum, fund, tax int

		err = tx.QueryRowContext(ctx, query, command.SalaryID).Scan(&salarySum, &fund, &tax)
		if err != nil {
			return err
		}
		taxed := int(float32(fund) * (float32(tax) / 100))
		fund -= taxed
		query = `UPDATE salaries SET status=$2 WHERE id = $1`
		if math.Abs(float64(fund-salarySum)) <= 1 {
			_, err = tx.ExecContext(ctx, query, command.SalaryID, serializable.SalaryPaid)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
