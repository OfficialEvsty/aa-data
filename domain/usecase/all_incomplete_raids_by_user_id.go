package usecase

import (
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

type IncompleteRaidDTO struct {
	UserID    uuid.UUID                 `json:"user_id"`
	RaidID    uuid.UUID                 `json:"raid_id"`
	PublishID uuid.UUID                 `json:"publish_id"`
	Status    serializable.Status       `json:"status"`
	S3Data    serializable.S3Screenshot `json:"s3_data"`
	CreatedAt time.Time                 `json:"created_at"`
	RaidAt    *time.Time                `json:"raid_at"`
	Events    []*EventDTO               `json:"events"`
	Validated bool                      `json:"validated"`
}
