package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	db2 "github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	repos2 "github.com/OfficialEvsty/aa-data/repos"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
	"time"
)

type AddLunarkAttendedToTenantCommand struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	LunarkID  uuid.UUID `json:"lunark_id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
}

type LunarkImporter struct {
	tx               *db.TxManager
	lunarkRepo       repos.ILunarkRepository
	tenantLunarkRepo junction_repos.ITenantLunarkRepository
}

func NewLunarkImporter(sqlExecutor *sql.DB) *LunarkImporter {
	return &LunarkImporter{
		tx:               db2.NewTxManager(sqlExecutor),
		lunarkRepo:       repos2.NewLunarkRepository(sqlExecutor),
		tenantLunarkRepo: junction_repos2.NewTenantLunarkRepository(sqlExecutor),
	}
}

func (li *LunarkImporter) Handle(ctx context.Context, cmd AddLunarkAttendedToTenantCommand) error {
	err := li.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		lunark := domain.Lunark{
			ID:        cmd.LunarkID,
			Name:      cmd.Name,
			StartDate: cmd.StartDate,
		}
		err := li.lunarkRepo.WithTx(tx).Add(ctx, lunark)
		if err != nil {
			return err
		}
		tl := domain.Journal{
			TenantID: cmd.TenantID,
			LunarkID: cmd.LunarkID,
		}
		err = li.tenantLunarkRepo.WithTx(tx).Add(ctx, tl)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
