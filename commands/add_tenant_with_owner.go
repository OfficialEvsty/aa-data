package commands

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
)

type AddTenantWithOwnerCommand struct {
	tenant domain.Tenant
}

type TenantConstructor struct {
	txManager      db.ITxExecutor
	tenantRepo     repos2.ITenantRepository
	tenantUserRepo junction_repos.ITenantUserRepository
}

func NewTenantConstructor(txManager db.ITxExecutor, tenantRepo repos2.ITenantRepository, tenantUserRepo junction_repos.ITenantUserRepository) *TenantConstructor {
	return &TenantConstructor{
		txManager:      txManager,
		tenantRepo:     tenantRepo,
		tenantUserRepo: tenantUserRepo,
	}
}

func (h *TenantConstructor) Handle(ctx context.Context, cmd AddTenantWithOwnerCommand) error {
	err := h.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		tenant, err := h.tenantRepo.WithTx(tx).Add(ctx, cmd.tenant)
		if err != nil {
			return err
		}
		err = h.tenantUserRepo.WithTx(tx).Add(ctx, tenant.ID, tenant.OwnerID)
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
