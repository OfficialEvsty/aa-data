package junction_repos

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
	"strings"
)

type RaidEventRepository struct {
	exec db.ISqlExecutor
}

func NewRaidEventRepository(exec db.ISqlExecutor) *RaidEventRepository {
	return &RaidEventRepository{exec}
}

func (r *RaidEventRepository) Add(ctx context.Context, re domain.RaidEvent) error {
	query := `INSERT INTO raid_events (raid_id, event_id) VALUES ($1, $2)`
	_, err := r.exec.ExecContext(ctx, query, re.RaidID, re.EventID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RaidEventRepository) AddMany(ctx context.Context, raidID uuid.UUID, eventIDs []int) error {
	valueStrings := make([]string, 0, len(eventIDs))
	valueArgs := make([]interface{}, 0, len(eventIDs)*2)

	for i, eventID := range eventIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, raidID, eventID)
	}

	stmt := fmt.Sprintf("INSERT INTO raid_events (raid_id, event_id) VALUES %s ON CONFLICT (raid_id, event_id) DO NOTHING", strings.Join(valueStrings, ","))
	_, err := r.exec.ExecContext(ctx, stmt, valueArgs...)
	return err
}

func (r *RaidEventRepository) Remove(ctx context.Context, rID uuid.UUID, eID int) error {
	query := `DELETE FROM raid_events WHERE raid_id = $1 AND event_id = $2`
	_, err := r.exec.ExecContext(ctx, query, rID, eID)
	if err != nil {
		return err
	}
	return nil
}
func (r *RaidEventRepository) All(ctx context.Context, rID uuid.UUID) ([]*domain.RaidEvent, error) {
	var raidEvents []*domain.RaidEvent
	query := `SELECT raid_id, event_id FROM raid_events WHERE raid_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, rID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var raidEvent domain.RaidEvent
		err := rows.Scan(&raidEvent.RaidID, &raidEvent.EventID)
		if err != nil {
			return nil, err
		}
		raidEvents = append(raidEvents, &raidEvent)
	}
	return raidEvents, nil
}
func (r *RaidEventRepository) WithTx(tx *sql.Tx) junction_repos2.IRaidEventRepository {
	return &RaidEventRepository{tx}
}
