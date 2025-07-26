package junction_repos

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

type FinishedPublishRepository struct {
	exec db.ISqlExecutor
}

func NewFinishedPublishRepository(exec db.ISqlExecutor) *FinishedPublishRepository {
	return &FinishedPublishRepository{exec}
}

func (r *FinishedPublishRepository) Add(ctx context.Context, publish domain.FinishedPublish) (*domain.FinishedPublish, error) {
	query := "INSERT INTO finished_publishes (publish_id, result) VALUES ($1, $2) RETURNING publish_id, result, finished_at"
	err := r.exec.QueryRowContext(
		ctx,
		query,
		publish.PublishID,
		publish.Result,
	).Scan(
		&publish.PublishID,
		&publish.Result,
		&publish.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &publish, nil
}

func (r *FinishedPublishRepository) Get(ctx context.Context, publishID uuid.UUID) (*domain.FinishedPublish, error) {
	var pub domain.FinishedPublish
	query := "SELECT * FROM finished_publishes WHERE publish_id = $1"
	err := r.exec.QueryRowContext(ctx, query, publishID).Scan(
		&pub.PublishID,
		&pub.Result,
		&pub.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &pub, nil
}

func (r *FinishedPublishRepository) Remove(ctx context.Context, publishID uuid.UUID) error {
	query := "DELETE FROM finished_publishes WHERE publish_id = $1"
	_, err := r.exec.ExecContext(ctx, query, publishID)
	return err
}
