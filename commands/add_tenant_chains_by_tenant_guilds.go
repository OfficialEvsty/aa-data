package commands

import (
	"context"
	"database/sql"
	"errors"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/queries"
	repos2 "github.com/OfficialEvsty/aa-data/repos"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
)

type AddTenantChainsCommand struct {
	GuildIDs []uuid.UUID
	TenantID uuid.UUID
}

type AddTenantChainsByTenantGuilds struct {
	txManager               *db.TxManager
	tenantChainRepo         junction_repos2.ITenantChainRepository
	chainRepo               repos.IChainRepository
	getActiveChainsByTenant *queries.GetActiveChainsByTenantQuery
	getAllMembersByGuildIDs *queries.GetAllMembersByGuildIDs
}

func NewAddTenantChainsByTenantGuilds(sql *sql.DB) *AddTenantChainsByTenantGuilds {
	return &AddTenantChainsByTenantGuilds{
		txManager:               db.NewTxManager(sql),
		tenantChainRepo:         junction_repos.NewTenantChainRepository(sql),
		chainRepo:               repos2.NewChainRepository(sql),
		getActiveChainsByTenant: queries.NewGetActiveChainsByTenantQuery(sql),
		getAllMembersByGuildIDs: queries.NewGetAllMembersByGuildIDs(sql),
	}
}

func (i *AddTenantChainsByTenantGuilds) Handle(ctx context.Context, cmd *AddTenantChainsCommand) error {
	return i.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		tenantChains, err := i.getActiveChainsByTenant.WithTx(tx).Handle(ctx, cmd.TenantID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		tenantNicknames, err := i.getAllMembersByGuildIDs.WithTx(tx).Handle(ctx, cmd.GuildIDs)
		if err != nil {
			return err
		}
		for _, tn := range tenantNicknames {
			exist := false
			for _, tc := range tenantChains {
				if tc.NicknameID == tn.NicknameID && tc.Active {
					exist = true
					break
				}
			}
			if !exist {
				chain := domain.NicknameChain{
					NicknameID:    tn.NicknameID,
					ChainID:       uuid.New(),
					ParentChainID: nil,
				}
				err = i.chainRepo.WithTx(tx).Add(ctx, chain)
				if err != nil {
					return err
				}
				tenantChain := domain.TenantChain{
					TenantID: cmd.TenantID,
					ChainID:  chain.ChainID,
				}
				err = i.tenantChainRepo.WithTx(tx).Add(ctx, tenantChain)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
