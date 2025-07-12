package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// IPublishRepository crud under screenshots data
type IPublishRepository interface {
	Add(ctx context.Context, screenshot domain.PublishedScreenshot) error
	Remove(ctx context.Context, screenshotID uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*domain.PublishedScreenshot, error)
	WithTx(*sql.Tx) IPublishRepository
}
