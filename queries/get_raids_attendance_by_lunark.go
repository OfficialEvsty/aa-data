package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetRaidsAttendanceByLunarkQuery struct {
	exec db.ISqlExecutor
}

func NewGetRaidsAttendanceByLunarkQuery(executor db.ISqlExecutor) *GetRaidsAttendanceByLunarkQuery {
	return &GetRaidsAttendanceByLunarkQuery{executor}
}

func (q *GetRaidsAttendanceByLunarkQuery) Handle(ctx context.Context, lunarkID uuid.UUID) (usecase.RaidsAttendanceByNicknameIDs, error) {
	query := `SELECT r.id, n.id
              FROM lunark_raids lr
              JOIN raids r ON r.id = lr.raid_id
              JOIN attendance a ON a.raid_id = lr.raid_id
              JOIN aa_nicknames n ON n.id = a.nickname_id
              WHERE lr.lunark_id = $1 AND r.is_deleted = FALSE`
	rows, err := q.exec.QueryContext(ctx, query, lunarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result = make(usecase.RaidsAttendanceByNicknameIDs)
	for rows.Next() {
		var nicknameID uuid.UUID
		var raidID uuid.UUID
		err = rows.Scan(
			&raidID,
			&nicknameID,
		)
		if err != nil {
			return nil, err
		}
		if result[raidID] == nil {
			result[raidID] = make(map[uuid.UUID]struct{})
		}

		result[raidID][nicknameID] = struct{}{}
	}
	return result, nil
}
