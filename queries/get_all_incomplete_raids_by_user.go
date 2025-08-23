package queries

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetAllIncompleteRaidsQuery struct {
	exec db.ISqlExecutor
}

func NewGetAllIncompleteRaidQuery(exec db.ISqlExecutor) *GetAllIncompleteRaidsQuery {
	return &GetAllIncompleteRaidsQuery{
		exec: exec,
	}
}

func (q *GetAllIncompleteRaidsQuery) Handle(ctx context.Context, userID uuid.UUID) ([]*usecase.IncompleteRaidDTO, error) {
	query := `
		SELECT 
			tp.user_id,
			r.id,
			tp.publish_id,
			r.status,
			r.created_at,
			r.raid_at,
			p.s3,
			EXISTS (
				SELECT 1 FROM raid_events re WHERE re.raid_id = r.id
			) AS has_events,
			re.event_id,
			te.name
		FROM raids AS r
		JOIN tenant_publishes AS tp ON tp.publish_id = r.publish_id
		JOIN publishes AS p ON p.id = r.publish_id
		LEFT JOIN raid_events AS re ON re.raid_id = r.id
		LEFT JOIN aa_template_events AS te ON re.event_id = te.id
		WHERE tp.user_id = $1
		  AND r.status <> 'resolved'
		  AND r.is_deleted = FALSE
		ORDER BY r.created_at DESC
	`

	rows, err := q.exec.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query raids: %w", err)
	}
	defer rows.Close()

	// словарь для агрегации рейдов
	raidMap := make(map[uuid.UUID]*usecase.IncompleteRaidDTO)

	for rows.Next() {
		var (
			raidID    uuid.UUID
			dto       usecase.IncompleteRaidDTO
			eventID   *int
			eventName *string
		)

		err := rows.Scan(
			&dto.UserID,
			&raidID,
			&dto.PublishID,
			&dto.Status,
			&dto.CreatedAt,
			&dto.RaidAt,
			&dto.S3Data,
			&dto.Validated,
			&eventID,
			&eventName,
		)
		if err != nil {
			return nil, fmt.Errorf("scan raid row: %w", err)
		}

		// если рейд ещё не в мапе — кладём туда
		raid, exists := raidMap[raidID]
		if !exists {
			dto.RaidID = raidID
			raidMap[raidID] = &dto
			raid = &dto
		}

		// если у строки есть событие — добавляем
		if eventID != nil {
			raid.Events = append(raid.Events, &usecase.EventDTO{
				TemplateID: *eventID,
				Name:       derefString(eventName),
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// собираем в слайс
	incompleteRaids := make([]*usecase.IncompleteRaidDTO, 0, len(raidMap))
	for _, raid := range raidMap {
		incompleteRaids = append(incompleteRaids, raid)
	}

	return incompleteRaids, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (q *GetAllIncompleteRaidsQuery) WithTx(tx *sql.Tx) *GetAllIncompleteRaidsQuery {
	return &GetAllIncompleteRaidsQuery{
		tx,
	}
}
