package commands

import (
	"context"
	"database/sql"
	db3 "github.com/OfficialEvsty/aa-data/db"
	db2 "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos2 "github.com/OfficialEvsty/aa-data/repos"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
)

type AddRequestToTenantCommand struct {
	TenantID uuid.UUID
	Request  domain.Request
}

type TenantRequestImporter struct {
	tx                db2.ITxExecutor
	requestRepo       repos.IRequestRepository
	tenantRequestRepo junction_repos.ITenantRequestRepository
}

func NewTenantRequestImporter(db *sql.DB) *TenantRequestImporter {
	return &TenantRequestImporter{
		tx:                db3.NewTxManager(db),
		requestRepo:       repos2.NewRequestRepository(db),
		tenantRequestRepo: junction_repos2.NewTenantRequestRepository(db),
	}
}

func (i *TenantRequestImporter) Handle(ctx context.Context, cmd *AddRequestToTenantCommand) error {
	return i.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := i.requestRepo.WithTx(tx).Add(ctx, cmd.Request)
		if err != nil {
			return err
		}
		var tr = domain.TenantRequest{
			TenantID:  cmd.TenantID,
			RequestID: cmd.Request.ID,
		}
		return i.tenantRequestRepo.WithTx(tx).Add(ctx, tr)
	})
}
