package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/OfficialEvsty/aa-data/errors"
	"github.com/OfficialEvsty/aa-data/repos"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

type AddDropAndNicknamesToRaidCommand struct {
	RaidID       uuid.UUID                `json:"raid_id"`
	Version      int64                    `json:"version"`
	NicknameIDs  []uuid.UUID              `json:"nickname_ids"`
	Attendance   int                      `json:"attendance"`
	DropItemList []*serializable.DropItem `json:"drop"`
}

type DropAndNicknamesImporter struct {
	txManager           *db.TxManager
	raidRepo            repos2.IRaidRepository
	addNicknamesCommand *AttendanceController
	addDropCommand      *DropItemCleanerAndImporter
}

func NewDropAndNicknamesImporter(sql *sql.DB) *DropAndNicknamesImporter {
	txManager := db.NewTxManager(sql)
	raidRepo := repos.NewRaidRepository(sql)
	addNicknamesCommand := NewAttendanceController(sql)
	addDropCommand := NewDropItemCleanerAndImporter(sql)
	return &DropAndNicknamesImporter{
		txManager:           txManager,
		raidRepo:            raidRepo,
		addNicknamesCommand: addNicknamesCommand,
		addDropCommand:      addDropCommand,
	}
}

func (i *DropAndNicknamesImporter) Handle(ctx context.Context, cmd *AddDropAndNicknamesToRaidCommand) error {
	return i.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		raid, err := i.raidRepo.WithTx(tx).GetById(ctx, cmd.RaidID)
		if err != nil {
			return err
		}
		if raid.Version != cmd.Version {
			return errors.ErrorRaidVersionMismatch
		}
		var raidNotBeenResolvedEarlier = raid.Status != serializable.StatusResolved &&
			(cmd.NicknameIDs == nil || len(cmd.NicknameIDs) == 0)
		if raidNotBeenResolvedEarlier {
			return errors.ErrorRaidPartialSavedRestricted
		}
		if len(cmd.NicknameIDs) > 0 {
			nCmd := &AddNicknamesToRaidCommand{
				RaidID:      cmd.RaidID,
				NicknameIDs: cmd.NicknameIDs,
				Attendance:  cmd.Attendance,
			}
			err = i.addNicknamesCommand.Handle(ctx, nCmd)
			if err != nil {
				return err
			}
		}
		if len(cmd.DropItemList) > 0 {
			dCmd := &ClearAndAddItemsAsRaidDropByRaidIDCommand{
				DropItemList: cmd.DropItemList,
				RaidID:       cmd.RaidID,
			}
			err = i.addDropCommand.Handle(ctx, dCmd)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
