package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/OfficialEvsty/aa-shared/golinq"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type GetAllCompletedRaidsByTenantID struct {
	exec db.ISqlExecutor
}

func NewGetAllCompletedRaidsByTenantID(exec db.ISqlExecutor) *GetAllCompletedRaidsByTenantID {
	return &GetAllCompletedRaidsByTenantID{exec: exec}
}

func (q *GetAllCompletedRaidsByTenantID) Handle(ctx context.Context, tenantID uuid.UUID) (*usecase.AllRaidsByTenantDTO, error) {
	getOpenedLunarkQuery := `SELECT l.id, l.name, l.start_date, l.end_date FROM tenant_lunark tl
							 JOIN lunark l ON l.id = tl.lunark_id           
							 WHERE tl.tenant_id = $1 AND l.opened = TRUE`
	var lunarkDTO usecase.LunarkDTO
	var lunarkID uuid.UUID
	err := q.exec.QueryRowContext(
		ctx,
		getOpenedLunarkQuery,
		tenantID,
	).Scan(
		&lunarkID,
		&lunarkDTO.Name,
		&lunarkDTO.StartDate,
		&lunarkDTO.EndDate,
	)

	getRaidsAndParticipantCountQuery := `SELECT r.id, tp.user_id, r.raid_at, COUNT(a.nickname_id) AS participants_count, r.attendance
 			  FROM raids AS r
 			  JOIN tenant_publishes AS tp ON r.publish_id = tp.publish_id
 			  JOIN attendance AS a ON r.id = a.raid_id
 			  JOIN lunark_raids AS lr ON r.id = lr.raid_id
 			  WHERE tp.tenant_id = $1 AND r.status = 'resolved' AND lr.lunark_id = $2
 			  GROUP BY r.id, tp.user_id, r.raid_at, r.attendance`

	rows, err := q.exec.QueryContext(ctx, getRaidsAndParticipantCountQuery, tenantID, lunarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	allRaidsDTO := usecase.AllRaidsByTenantDTO{
		Raids:  make(map[uuid.UUID]*usecase.RaidDTO),
		Events: make(map[int]*usecase.EventDTO),
		Lunark: &lunarkDTO,
	}
	var attendanceSum uint
	for rows.Next() {
		var r usecase.RaidDTO
		err = rows.Scan(
			&r.ID,
			&r.ProviderID,
			&r.OccurredAt,
			&r.ParticipantCount,
			&r.Attendance,
		)
		if err != nil {
			return nil, err
		}
		attendanceSum = attendanceSum + r.Attendance
		allRaidsDTO.Raids[r.ID] = &r
	}
	if len(allRaidsDTO.Raids) != 0 {
		allRaidsDTO.Attendance = attendanceSum / uint(len(allRaidsDTO.Raids))
	}
	getRaidEventsQuery := `SELECT r.id, re.event_id, te.name
            			   FROM raid_events AS re
            			   JOIN raids AS r ON r.id = re.raid_id
            			   JOIN aa_template_events AS te ON re.event_id = te.id
            			   WHERE re.raid_id = ANY($1) AND r.is_deleted = FALSE`
	var raids []*usecase.RaidDTO
	for _, val := range allRaidsDTO.Raids {
		raids = append(raids, val)
	}
	raidIDs := golinq.Map(raids, func(r *usecase.RaidDTO) uuid.UUID { return r.ID })
	rows, err = q.exec.QueryContext(ctx, getRaidEventsQuery, pq.Array(raidIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e usecase.EventDTO
		var rID uuid.UUID
		err = rows.Scan(
			&rID,
			&e.TemplateID,
			&e.Name,
		)
		if err != nil {
			return nil, err
		}
		raid := allRaidsDTO.Raids[rID]
		raid.EventIDs = append(allRaidsDTO.Raids[rID].EventIDs, e.TemplateID)
		allRaidsDTO.Events[e.TemplateID] = &e
	}
	return &allRaidsDTO, nil
}

func (q *GetAllCompletedRaidsByTenantID) WithTx(tx *sql.Tx) *GetAllCompletedRaidsByTenantID {
	return &GetAllCompletedRaidsByTenantID{
		exec: tx,
	}
}
