package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
)

type IRawOcrRepository interface {
	Add(context.Context, uuid.UUID, serializable.OCRData) error
	WithTx(*sql.Tx) IRawOcrRepository
}
