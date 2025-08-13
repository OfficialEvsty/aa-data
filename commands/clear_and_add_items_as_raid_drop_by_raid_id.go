package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type ClearAndAddItemsAsRaidDropByRaidIDCommand struct {
	DropItemList []*serializable.DropItem
	RaidID       uuid.UUID
}

type DropItemCleanerAndImporter struct {
	txManager    *db.TxManager
	raidItemRepo junction_repos.IRaidItemRepository
}

func NewDropItemCleanerAndImporter(manager *db.TxManager, raidItemRepo junction_repos.IRaidItemRepository) *DropItemCleanerAndImporter {
	return &DropItemCleanerAndImporter{
		txManager:    manager,
		raidItemRepo: raidItemRepo,
	}
}

func (ci *DropItemCleanerAndImporter) Handle(ctx context.Context, cmd *ClearAndAddItemsAsRaidDropByRaidIDCommand) error {
	err := ci.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := ci.raidItemRepo.WithTx(tx).RemoveItemsByRaidID(ctx, cmd.RaidID)
		if err != nil {
			return err
		}
		err = ci.raidItemRepo.WithTx(tx).AddOrUpdateItems(ctx, cmd.RaidID, cmd.DropItemList)
		return err
	})
	return err
}
