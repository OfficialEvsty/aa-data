package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetRaidItemsQuery struct {
	exec db.ISqlExecutor
}

func NewGetRaidItemsQuery(exec db.ISqlExecutor) *GetRaidItemsQuery {
	return &GetRaidItemsQuery{exec}
}

// deprecated
func (q *GetRaidItemsQuery) Handle(ctx context.Context, raidID uuid.UUID) ([]*usecase.RaidItemDTO, error) {
	query := `SELECT i.id, i.name, i.img_url, ri.rate
              FROM raid_items ri
              JOIN aa_items i ON ri.item_id = i.id
              JOIN raids r ON ri.raid_id = r.id
              WHERE ri.raid_id = $1 AND r.is_deleted = FALSE`
	rows, err := q.exec.QueryContext(ctx, query, raidID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*usecase.RaidItemDTO
	for rows.Next() {
		var item usecase.RaidItemDTO
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.ImageURL,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}
