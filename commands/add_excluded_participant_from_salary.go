package commands

import (
	"context"
	"database/sql"
	"errors"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	errors2 "github.com/OfficialEvsty/aa-data/errors"
	"github.com/OfficialEvsty/aa-data/repos"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
)

type AddExcludedParticipantFromSalaryCommand struct {
	ChainID        uuid.UUID
	SalaryID       uuid.UUID
	Reason         string
	CurrentVersion int
}

type RemoveExcludedParticipantFromSalaryCommand struct {
	ChainID        uuid.UUID
	SalaryID       uuid.UUID
	CurrentVersion int
}

type SalaryExcludingParticipantManager struct {
	txManager               *db.TxManager
	excludedParticipantRepo junction_repos.IExcludedParticipantRepository
	salaryRepo              repos2.ISalaryRepository
}

func NewSalaryExcludingParticipantManager(sql *sql.DB) *SalaryExcludingParticipantManager {
	return &SalaryExcludingParticipantManager{
		txManager:               db.NewTxManager(sql),
		excludedParticipantRepo: junction_repos2.NewExcludedParticipantRepository(sql),
		salaryRepo:              repos.NewSalaryRepository(sql),
	}
}

func (m *SalaryExcludingParticipantManager) Handle(
	ctx context.Context,
	command *AddExcludedParticipantFromSalaryCommand,
) error {
	return m.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		exP := usecase.ExcludedParticipant{
			ChainID:  command.ChainID,
			SalaryID: command.SalaryID,
			Reason:   command.Reason,
		}
		err := m.excludedParticipantRepo.WithTx(tx).Add(ctx, exP)
		if err != nil {
			return err
		}
		salary, err := m.salaryRepo.WithTx(tx).GetByID(ctx, command.SalaryID)
		if err != nil {
			return err
		}
		if salary.Version != command.CurrentVersion {
			return errors2.ErrorSalaryVersionMismatch
		}
		err = m.salaryRepo.WithTx(tx).Update(ctx, *salary)
		if err != nil {
			return err
		}
		return nil
	})
}

func (m *SalaryExcludingParticipantManager) RemoveExcludedParticipant(
	ctx context.Context,
	cmd *RemoveExcludedParticipantFromSalaryCommand,
) error {
	return m.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := m.excludedParticipantRepo.WithTx(tx).Remove(ctx, cmd.ChainID, cmd.SalaryID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		salary, err := m.salaryRepo.WithTx(tx).GetByID(ctx, cmd.SalaryID)
		if err != nil {
			return err
		}
		if salary.Version != cmd.CurrentVersion {
			return errors2.ErrorSalaryVersionMismatch
		}
		err = m.salaryRepo.WithTx(tx).Update(ctx, *salary)
		if err != nil {
			return err
		}
		return nil
	})
}
