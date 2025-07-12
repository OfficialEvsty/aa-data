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

// UserPublishedScreenshot it's PublishedScreenshot bounds to specified user id
type UserPublishedScreenshot struct {
	PublishedID uuid.UUID `json:"published_id"`
	UserID      uuid.UUID `json:"user_id"`
	PublishedAt time.Time `json:"published_at"`
}
