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
		return nil
	})
}
