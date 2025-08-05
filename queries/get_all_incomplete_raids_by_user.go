package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

type IncompleteRaidDTO struct {
	UserID    uuid.UUID           `json:"user_id"`
	RaidID    uuid.UUID           `json:"raid_id"`
	PublishID uuid.UUID           `json:"publish_id"`
	Status    serializable.Status `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	RaidAt    *time.Time          `json:"raid_at"`
}

type GetAllIncompleteRaidsQuery struct {
	exec db.ISqlExecutor
}

func NewGetAllIncompleteRaidQuery(exec db.ISqlExecutor) *GetAllIncompleteRaidsQuery {
	return &GetAllIncompleteRaidsQuery{
		exec: exec,
	}
}

func (q *GetAllIncompleteRaidsQuery) Handle(ctx context.Context, userID uuid.UUID) ([]*IncompleteRaidDTO, error) {
	var incompleteRaids []*IncompleteRaidDTO
	query := `SELECT tp.user_id, r.id, tp.publish_id, r.status, r.created_at, r.raid_at
			  FROM raids AS r 
			  JOIN tenant_publishes AS tp ON tp.publish_id = r.publish_id
			  WHERE tp.user_id = $1`
	rows, err := q.exec.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var incompleteRaidDTO IncompleteRaidDTO
		err = rows.Scan(
			&incompleteRaidDTO.UserID,
			&incompleteRaidDTO.RaidID,
			&incompleteRaidDTO.PublishID,
			&incompleteRaidDTO.Status,
			&incompleteRaidDTO.CreatedAt,
			&incompleteRaidDTO.RaidAt,
		)
		if err != nil {
			return nil, err
		}
		incompleteRaids = append(incompleteRaids, &incompleteRaidDTO)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return incompleteRaids, nil
}

func (q *GetAllIncompleteRaidsQuery) WithTx(tx *sql.Tx) *GetAllIncompleteRaidsQuery {
	return &GetAllIncompleteRaidsQuery{
		tx,
	}
}
