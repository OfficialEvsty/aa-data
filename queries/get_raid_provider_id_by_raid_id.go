package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/google/uuid"
)

type GetRaidProviderIdByRaidIdQuery struct {
	exec db.ISqlExecutor
}

func NewGetRaidProviderIdByRaidIdQuery(db db.ISqlExecutor) *GetRaidProviderIdByRaidIdQuery {
	return &GetRaidProviderIdByRaidIdQuery{db}
}

func (q *GetRaidProviderIdByRaidIdQuery) Handle(ctx context.Context, raidID uuid.UUID) (uuid.UUID, error) {
	var providerID uuid.UUID
	query := `SELECT tp.user_id
              FROM tenant_publishes tp
              JOIN raids r ON r.publish_id = tp.publish_id
              WHERE r.raid_id = $1`
	err := q.exec.QueryRowContext(ctx, query, raidID).Scan(&providerID)
	if err != nil {
		return uuid.Nil, err
	}
	return providerID, nil
}

func (q *GetRaidProviderIdByRaidIdQuery) WithTx(tx *sql.Tx) *GetRaidProviderIdByRaidIdQuery {
	return &GetRaidProviderIdByRaidIdQuery{tx}
}
