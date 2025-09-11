package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	repos2 "github.com/OfficialEvsty/aa-data/repos"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

type AttachNewChainByOldChainIDCommand struct {
	ChildChainID  uuid.UUID
	ParentChainID uuid.UUID
}

type DetachChainFromParentCommand struct {
	ParentChainID uuid.UUID
}

type AttachManager struct {
	txManager *db.TxManager
	chainRepo repos.IChainRepository
}

func NewAttachManager(sql *sql.DB) *AttachManager {
	return &AttachManager{
		txManager: db.NewTxManager(sql),
		chainRepo: repos2.NewChainRepository(sql),
	}
}

func (ci *AttachManager) AttachChain(ctx context.Context, cmd *AttachNewChainByOldChainIDCommand) error {
	return ci.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := ci.chainRepo.WithTx(tx).AttachChain(ctx, cmd.ParentChainID, cmd.ChildChainID)
		if err != nil {
			return err
		}
		err = ci.chainRepo.WithTx(tx).UpdateStatus(ctx, cmd.ParentChainID, false)
		if err != nil {
			return err
		}
		return nil
	})
}

func (ci *AttachManager) DetachChain(ctx context.Context, cmd *DetachChainFromParentCommand) error {
	return ci.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := ci.chainRepo.WithTx(tx).DetachChain(ctx, cmd.ParentChainID)
		if err != nil {
			return err
		}
		err = ci.chainRepo.WithTx(tx).UpdateStatus(ctx, cmd.ParentChainID, true)
		if err != nil {
			return err
		}
		return nil
	})
}
