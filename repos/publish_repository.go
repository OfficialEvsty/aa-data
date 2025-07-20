package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

// PublishRepository implementation of IPublishRepository
type PublishRepository struct {
	exec db.ISqlExecutor
}

func NewPublishRepository(exec db.ISqlExecutor) *PublishRepository {
	return &PublishRepository{exec}
}

// Add implements Add method
func (r *PublishRepository) Add(ctx context.Context, s domain.PublishedScreenshot) error {
	query := `INSERT INTO publishes (id, s3) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING;`
	_, err := r.exec.ExecContext(ctx, query, s.ID, s.S3Data)
	if err != nil {
		return err
	}
	return nil
}

func (r *PublishRepository) Remove(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM publishes WHERE id = $1;`
	_, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PublishRepository) Get(ctx context.Context, id uuid.UUID) (*domain.PublishedScreenshot, error) {
	var s domain.PublishedScreenshot
	query := `SELECT id, s3 FROM publishes WHERE id = $1;`
	row := r.exec.QueryRowContext(ctx, query, id)
	err := row.Scan(&s.ID, &s.S3Data)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PublishRepository) WithTx(tx *sql.Tx) repos.IPublishRepository {
	return &PublishRepository{tx}
}
