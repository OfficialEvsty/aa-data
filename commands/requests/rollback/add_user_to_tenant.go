package rollback

import (
	"context"
	"database/sql"
	"encoding/json"
	db2 "github.com/OfficialEvsty/aa-data/db"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	repos2 "github.com/OfficialEvsty/aa-data/repos"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
)

type AddUserToTenantPayload struct {
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
}

type AddUserToTenantRequest struct {
	cmd            AddUserToTenantPayload
	tx             db.ITxExecutor
	tenantUserRepo junction_repos.ITenantUserRepository
	requestRepo    repos.IRequestRepository
}

func NewAddUserToTenantRequest(sql *sql.DB) *AddUserToTenantRequest {

	return &AddUserToTenantRequest{
		tx:             db2.NewTxManager(sql),
		tenantUserRepo: junction_repos2.NewTenantUserRepository(sql),
		requestRepo:    repos2.NewRequestRepository(sql),
	}
}

func (r *AddUserToTenantRequest) Execute(ctx context.Context, payload []byte) error {
	err := json.Unmarshal(payload, &r.cmd)
	if err != nil {
		return err
	}
	return r.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		exists, err := r.tenantUserRepo.WithTx(tx).CheckUser(ctx, r.cmd.TenantID, r.cmd.UserID)
		if err != nil {
			return err
		}
		if !exists {
			err = r.tenantUserRepo.WithTx(tx).Add(ctx, r.cmd.TenantID, r.cmd.UserID)
			if err != nil {
				return err
			}
		}
		err = r.requestRepo.WithTx(tx).Accept(ctx, r.cmd.TenantID, r.cmd.UserID)
		return err
	})
}

func (r *AddUserToTenantRequest) Rollback(ctx context.Context, payload []byte) error {
	err := json.Unmarshal(payload, &r.cmd)
	if err != nil {
		return err
	}

	return r.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		exists, err := r.tenantUserRepo.WithTx(tx).CheckUser(ctx, r.cmd.TenantID, r.cmd.UserID)
		if err != nil {
			return err
		}
		if exists {
			err = r.tenantUserRepo.WithTx(tx).Remove(ctx, r.cmd.TenantID, r.cmd.UserID)
			if err != nil {
				return err
			}
		}
		return r.requestRepo.WithTx(tx).Remove(ctx, r.cmd.UserID)
	})
}

func (r *AddUserToTenantRequest) Type() serializable.RequestType {
	return serializable.AddUserToTenantRequest
}
