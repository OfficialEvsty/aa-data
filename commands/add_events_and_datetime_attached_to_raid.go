package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/repos"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/google/uuid"
	"time"
)

type AddEventsAndDatetimeAttachedToRaid struct {
	RaidID   uuid.UUID `json:"raid"`
	EventIDs []int     `json:"event_ids"`
	DateTime time.Time `json:"datetime"`
}

type EventAndDatetimeImporter struct {
	txManager     *db.TxManager
	raidEventRepo junction_repos.IRaidEventRepository
	raidRepo      repos2.IRaidRepository
}

func NewEventAndDatetimeImporter(sql *sql.DB) *EventAndDatetimeImporter {
	return &EventAndDatetimeImporter{
		txManager:     db.NewTxManager(sql),
		raidRepo:      repos.NewRaidRepository(sql),
		raidEventRepo: junction_repos2.NewRaidEventRepository(sql),
	}
}

func (i *EventAndDatetimeImporter) Handle(ctx context.Context, cmd *AddEventsAndDatetimeAttachedToRaid) error {
	return i.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := i.raidRepo.UpdateTiming(ctx, cmd.RaidID, cmd.DateTime)
		if err != nil {
			return err
		}
		err = i.raidEventRepo.AddMany(ctx, cmd.RaidID, cmd.EventIDs)
		if err != nil {
			return err
		}
		return nil
	})
}
