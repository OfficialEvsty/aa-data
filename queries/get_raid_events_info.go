package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetRaidEventsInfoQuery struct {
	exec db.ISqlExecutor
}

func NewGetRaidEventsInfoQuery(exec db.ISqlExecutor) *GetRaidEventsInfoQuery {
	return &GetRaidEventsInfoQuery{exec}
}

func (q *GetRaidEventsInfoQuery) Handle(ctx context.Context, raidID uuid.UUID) ([]*usecase.RaidEventDTO, error) {
	query := `SELECT re.raid_id, re.event_id, e.name
              FROM raid_events re
              JOIN aa_template_events e ON re.event_id = e.id
              JOIN raids r ON re.raid_id = r.id
              WHERE r.id = $1 AND r.is_deleted = FALSE`
	rows, err := q.exec.QueryContext(ctx, query, raidID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*usecase.RaidEventDTO
	for rows.Next() {
		var event usecase.RaidEventDTO
		if err = rows.Scan(&event.RaidID, &event.EventID, &event.Name); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
