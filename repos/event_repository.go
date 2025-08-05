package repos

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
)

type EventRepository struct {
	exec db.ISqlExecutor
}

func NewEventRepository(exec db.ISqlExecutor) *EventRepository {
	return &EventRepository{exec}
}

func (r *EventRepository) Add(ctx context.Context, event domain.Event) (*domain.Event, error) {
	query := `INSERT INTO aa_events (template_id) VALUES ($1) RETURNING id, template_id, occurred_at;`
	err := r.exec.QueryRowContext(
		ctx,
		query,
		event.TemplateID,
	).Scan(&event.ID, &event.TemplateID, &event.OccurredAt)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id uint64) (*domain.Event, error) {
	var event domain.Event
	query := `SELECT * FROM aa_events WHERE id = $1;`
	err := r.exec.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(&event.ID, &event.TemplateID, &event.OccurredAt)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Remove(ctx context.Context, id uint64) error {
	query := `DELETE FROM aa_events WHERE id = $1;`
	_, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
