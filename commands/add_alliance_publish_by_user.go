package commands

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type AddTenantPublishByUser struct {
	PublishID uuid.UUID
	TenantID  uuid.UUID
	UserID    uuid.UUID
	S3Data    serializable.S3Screenshot
}

// Publisher make general transaction to tenant, user and publish tables
type Publisher struct {
	tx                db.ITxExecutor
	publishRepo       repos.IPublishRepository
	tenantPublishRepo junction_repos.ITenantPublishRepository
}

func NewPublisher(
	tx db.ITxExecutor,
	publishRepo repos.IPublishRepository,
	tenantPublishRepo junction_repos.ITenantPublishRepository,
) *Publisher {
	return &Publisher{
		tx:                tx,
		publishRepo:       publishRepo,
		tenantPublishRepo: tenantPublishRepo,
	}
}

func (si *Publisher) Handle(ctx context.Context, cmd AddTenantPublishByUser) error {
	err := si.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := si.publishRepo.WithTx(tx).Add(ctx, domain.PublishedScreenshot{ID: cmd.PublishID, S3Data: cmd.S3Data})
		if err != nil {
			return err
		}
		_, err = si.tenantPublishRepo.WithTx(tx).Add(ctx, domain.TenantPublish{TenantID: cmd.TenantID, PublishID: cmd.PublishID, UserID: cmd.UserID})
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
