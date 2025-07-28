package commands

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type AddTenantPublishByUserCommand struct {
	PublishID uuid.UUID                 `json:"publish_id"`
	S3Data    serializable.S3Screenshot `json:"s3_data"`
	TenantID  uuid.UUID                 `json:"tenant_id"`
	UserID    uuid.UUID                 `json:"user_id"`
}

type TenantPublisher struct {
	txManager         db.TxManager
	tenantPublishRepo junction_repos.ITenantPublishRepository
	publishRepo       repos.IPublishRepository
}

func NewTenantPublisher(
	txManager db.TxManager,
	tenantPublishRepo junction_repos.ITenantPublishRepository,
	publishRepo repos.IPublishRepository,
) *TenantPublisher {
	return &TenantPublisher{
		txManager:         txManager,
		tenantPublishRepo: tenantPublishRepo,
		publishRepo:       publishRepo,
	}
}

func (tp *TenantPublisher) Handle(ctx context.Context, cmd *AddTenantPublishByUserCommand) error {
	err := tp.txManager.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		publish := domain.PublishedScreenshot{
			ID:     cmd.PublishID,
			S3Data: cmd.S3Data,
		}
		err := tp.publishRepo.WithTx(tx).Add(ctx, publish)
		if err != nil {
			return err
		}
		tenantPublish := domain.TenantPublish{
			UserID:    cmd.UserID,
			TenantID:  cmd.TenantID,
			PublishID: cmd.PublishID,
		}
		_, err = tp.tenantPublishRepo.WithTx(tx).Add(ctx, tenantPublish)
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
