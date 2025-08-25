package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

type RawOcrRepository struct {
	exec db.ISqlExecutor
}

func NewRawOcrRepository(exec db.ISqlExecutor) *RawOcrRepository {
	return &RawOcrRepository{exec}
}

func (r *RawOcrRepository) Add(ctx context.Context, pubID uuid.UUID, result serializable.OCRData) error {
	query := `INSERT INTO raw_ocr_data (publish_id, raw) VALUES ($1, $2)`
	_, err := r.exec.ExecContext(ctx, query, pubID, result)
	if err != nil {
		return err
	}
	return nil
}

func (r *RawOcrRepository) WithTx(tx *sql.Tx) repos.IRawOcrRepository {
	return &RawOcrRepository{
		tx,
	}
}
