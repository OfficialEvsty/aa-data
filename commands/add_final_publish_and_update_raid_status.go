package commands

import (
	"context"
	"database/sql"
	"errors"
	db2 "github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type AddFinalPublishAndUpdateRaidStatus struct {
	PublishID uuid.UUID
	Data      serializable.NicknameResultWithConflicts
	IsError   bool
}

type FinalPublishImporter struct {
	txManager    db2.TxManager
	finalPubRepo junction_repos2.IFinishedPublish
	raidRepo     repos.IRaidRepository
}

func (fpi *FinalPublishImporter) Handle(ctx context.Context, cmd *AddFinalPublishAndUpdateRaidStatus) error {
	err := fpi.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		status := serializable.StatusRecognizeError
		if !cmd.IsError {
			finPub := domain.FinishedPublish{
				PublishID: cmd.PublishID,
				Result:    cmd.Data,
			}
			_, err := fpi.finalPubRepo.WithTx(tx).Add(ctx, finPub)
			if err != nil {
				return err
			}
			status = serializable.StatusRecognized
		}
		raid, err := fpi.raidRepo.WithTx(tx).GetByPublishID(ctx, cmd.PublishID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			return nil
		}
		err = fpi.raidRepo.WithTx(tx).UpdateStatus(ctx, raid.ID, status)
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
