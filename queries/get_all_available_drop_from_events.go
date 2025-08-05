package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/lib/pq"
)

type DropItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Tier     int    `json:"tier"`
	ImageURL string `json:"image_url"`
	TierURL  string `json:"tier_url"`
	Quantity string `json:"quantity"`
}

// Drops groupped with bosses and events
type EventBossLootDTO map[int]map[int64][]*DropItem

type GetAllAvailableDropFromEvents struct {
	exec db.ISqlExecutor
}

func NewGetAllAvailableDropFromEvents(exec db.ISqlExecutor) *GetAllAvailableDropFromEvents {
	return &GetAllAvailableDropFromEvents{
		exec: exec,
	}
}

func (q *GetAllAvailableDropFromEvents) Handle(ctx context.Context, templateEventIDs []int) (EventBossLootDTO, error) {
	var dto EventBossLootDTO
	dto = make(map[int]map[int64][]*DropItem)
	query := `SELECT eb.event_template_id, b.id, i.id, i.name, i.tier, i.img_url, i.img_grade_url, d.elem->>'rate'
    		  FROM aa_event_bosses AS eb
    		  JOIN aa_bosses AS b ON eb.boss_id = b.id
    		  LEFT JOIN LATERAL (
				  SELECT elem FROM jsonb_array_elements(b.drop) AS elem
				  WHERE jsonb_typeof(b.drop) = 'array'
				) AS d ON TRUE
    		  JOIN aa_items AS i ON (d.elem->>'item_id')::int = i.id
    		  WHERE eb.event_template_id = ANY($1);`
	rows, err := q.exec.QueryContext(ctx, query, pq.Array(templateEventIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var drop DropItem
		var eID int
		var bID int64
		err = rows.Scan(
			&eID,
			&bID,
			&drop.ID,
			&drop.Name,
			&drop.Tier,
			&drop.ImageURL,
			&drop.TierURL,
			&drop.Quantity,
		)
		if err != nil {
			return nil, err
		}
		if dto[eID] == nil {
			dto[eID] = make(map[int64][]*DropItem)
		}
		dto[eID][bID] = append(dto[eID][bID], &drop)
	}
	return dto, nil
}

func (q *GetAllAvailableDropFromEvents) WithTx(tx *sql.Tx) *GetAllAvailableDropFromEvents {
	return &GetAllAvailableDropFromEvents{
		tx,
	}
}
