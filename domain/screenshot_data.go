package domain

import (
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

// PublishedScreenshot screenshot uploaded on s3 storage and enqueues in rabbitmq to receive player's nicknames
type PublishedScreenshot struct {
	ID     uuid.UUID                 `json:"id"`
	S3Data serializable.S3Screenshot `json:"s3"`
}

type TenantPublish struct {
	UserID      uuid.UUID `json:"user_id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	PublishID   uuid.UUID `json:"publish_id"`
	PublishedAt time.Time `json:"published_at"`
}
