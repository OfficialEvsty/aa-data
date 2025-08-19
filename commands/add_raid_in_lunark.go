package commands

import (
	"context"
	"database/sql"
	"errors"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/OfficialEvsty/aa-data/repos"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
	"time"
)

type AddRaidInLunarkCommand struct {
	LunarkID   uuid.UUID           `json:"lunark_id"`
	RaidID     uuid.UUID           `json:"raid_id"`
	Status     serializable.Status `json:"status"`
	RaidAt     *time.Time          `json:"raid_at"`
	Attendance int                 `json:"attendance"`
}

type RaidImporter struct {
	tx             *db.TxManager
	lunarkRaidRepo junction_repos.ILunarkRaidRepository
	lunarkRepo     repos2.ILunarkRepository
	raidRepo       repos2.IRaidRepository
}

func NewRaidInLunarkCommand(exec *sql.DB) *RaidImporter {
	return &RaidImporter{
		tx:             db.NewTxManager(exec),
		lunarkRaidRepo: junction_repos2.NewLunarkRaidRepository(exec),
		lunarkRepo:     repos.NewLunarkRepository(exec),
		raidRepo:       repos.NewRaidRepository(exec),
	}
}

func (ri *RaidImporter) Handle(ctx context.Context, cmd AddRaidInLunarkCommand) error {
	err := ri.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		lr := domain.LunarkRaid{
			LunarkID: cmd.LunarkID,
			RaidID:   cmd.RaidID,
		}
		err := ri.lunarkRaidRepo.WithTx(tx).Add(ctx, lr)
		if err != nil {
			return err
		}
		if cmd.RaidAt == nil {
			return errors.New("raid_at is required")
		}
		RaidAt := *cmd.RaidAt
		err = ri.raidRepo.WithTx(tx).UpdateStatus(ctx, cmd.RaidID, cmd.Status)
		if err != nil {
			return err
		}
		err = ri.raidRepo.WithTx(tx).UpdateVersion(ctx, cmd.RaidID)
		if err != nil {
			return err
		}
		err = ri.lunarkRepo.WithTx(tx).UpdateEndDate(ctx, cmd.LunarkID, RaidAt)
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
