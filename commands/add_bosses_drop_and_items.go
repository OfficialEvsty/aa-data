package commands

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
)

// AddBossesDropAndItemsCommand adds bosses and their drop info in db
type AddBossesDropAndItemsCommand struct {
	Bosses []*domain.AABoss         `json:"bosses"`
	Items  []*domain.AAItemTemplate `json:"items"`
}

type BossesImporter struct {
	tx       db.ITxExecutor
	bossRepo repos2.IBossesRepository
	itemRepo repos2.IItemRepository
}

func NewBossesImporter(tx db.ITxExecutor,
	bossRepo repos2.IBossesRepository,
	itemRepo repos2.IItemRepository,
) *BossesImporter {
	return &BossesImporter{
		tx:       tx,
		bossRepo: bossRepo,
		itemRepo: itemRepo,
	}
}

// Handle AddBossesDropAndItemsCommand
func (si *BossesImporter) Handle(ctx context.Context, cmd AddBossesDropAndItemsCommand) error {
	err := si.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		for _, boss := range cmd.Bosses {
			_, err := si.bossRepo.WithTx(tx).Add(ctx, *boss)
			if err != nil {
				return fmt.Errorf("error adding a boss record in table: %v", err)
			}
		}
		for _, item := range cmd.Items {
			_, err := si.itemRepo.WithTx(tx).Add(ctx, *item)
			if err != nil {
				return fmt.Errorf("error adding an item record in table: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
