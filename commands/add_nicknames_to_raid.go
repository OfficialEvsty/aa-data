package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/repos"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
)

type AddNicknamesToRaidCommand struct {
	RaidID      uuid.UUID   `json:"raid_id"`
	NicknameIDs []uuid.UUID `json:"nickname_ids"`
	Attendance  int         `json:"attendance"`
}

type AttendanceController struct {
	tx               *db.TxManager
	raidRepo         repos2.IRaidRepository
	raidNicknameRepo junction_repos2.IRaidNicknameRepository
}

func NewAttendanceController(sql *sql.DB) *AttendanceController {
	return &AttendanceController{
		tx:               db.NewTxManager(sql),
		raidRepo:         repos.NewRaidRepository(sql),
		raidNicknameRepo: junction_repos.NewRaidNicknameRepository(sql),
	}
}

func (c *AttendanceController) Handle(ctx context.Context, cmd *AddNicknamesToRaidCommand) error {
	err := c.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := c.raidRepo.WithTx(tx).UpdateAttendance(ctx, cmd.RaidID, cmd.Attendance)
		if err != nil {
			return err
		}
		err = c.raidNicknameRepo.WithTx(tx).ClearNicknamesByRaidID(ctx, cmd.RaidID)
		if err != nil {
			return err
		}
		err = c.raidNicknameRepo.WithTx(tx).AddNicknames(ctx, cmd.RaidID, cmd.NicknameIDs)
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
