package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/OfficialEvsty/aa-data/repos"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
)

type AddSalaryContextByLunarkIDCommand struct {
	UserID        uuid.UUID
	LunarkID      uuid.UUID
	Fond          int64
	MinAttendance int
	Tax           int
	Status        serializable.SalaryStatus
}

type SalaryContextImporter struct {
	txManager        *db.TxManager
	lunarkSalaryRepo junction_repos.ILunarkSalaryRepository
	salaryRepo       repos2.ISalaryRepository
}

func NewSalaryContextImporter(sql *sql.DB) *SalaryContextImporter {
	return &SalaryContextImporter{
		txManager:        db.NewTxManager(sql),
		lunarkSalaryRepo: junction_repos2.NewLunarkSalaryRepository(sql),
		salaryRepo:       repos.NewSalaryRepository(sql),
	}
}

func (i *SalaryContextImporter) Handle(ctx context.Context, cmd *AddSalaryContextByLunarkIDCommand) error {
	return i.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		salary := domain.Salary{
			ID:            uuid.New(),
			Fond:          cmd.Fond,
			MinAttendance: cmd.MinAttendance,
			Tax:           cmd.Tax,
			SubmittedBy:   cmd.UserID,
			Status:        serializable.SalaryInProcessing,
		}
		err := i.salaryRepo.WithTx(tx).Add(ctx, salary)
		if err != nil {
			return err
		}
		lunarkSalary := domain.LunarkSalary{
			SalaryID: salary.ID,
			LunarkID: cmd.LunarkID,
		}
		err = i.lunarkSalaryRepo.WithTx(tx).Add(ctx, lunarkSalary)
		if err != nil {
			return err
		}
		return nil
	})
}
