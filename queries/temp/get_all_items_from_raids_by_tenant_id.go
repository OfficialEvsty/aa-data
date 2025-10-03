package temp

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetAllItemsFromRaidsByTenantIdQuery struct {
	exec db.ISqlExecutor
}

func NewGetAllItemsFromRaidsByTenantIdQuery(sql db.ISqlExecutor) *GetAllItemsFromRaidsByTenantIdQuery {
	return &GetAllItemsFromRaidsByTenantIdQuery{sql}
}

func (q *GetAllItemsFromRaidsByTenantIdQuery) Handle(
	ctx context.Context,
	tenantID uuid.UUID,
) (*usecase.TreasuryItemList, error) {
	query := `SELECT i.id, i.name, i.img_url, ri.rate
              FROM tenant_lunark tl
              JOIN lunark_raids lr ON lr.lunark_id = tl.lunark_id
              JOIN raid_items ri ON ri.raid_id = lr.raid_id
			  JOIN aa_items i ON i.id = ri.item_id
              JOIN raids r ON r.id = lr.raid_id
			  WHERE tl.tenant_id = $1 AND r.is_deleted = FALSE`
	rows, err := q.exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items map[uint64]*usecase.TreasuryItem = make(map[uint64]*usecase.TreasuryItem)
	var itemList usecase.TreasuryItemList
	for rows.Next() {
		var item usecase.TreasuryItem
		err = rows.Scan(
			&item.ItemID,
			&item.Name,
			&item.IconURL,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}
		// Эссенция ярости
		if item.ItemID == 48404 {
			itemList.FierceEssenceQuantity += item.Quantity
			continue
		} else if item.ItemID == 54829 {
			itemList.TrophyEssenceQuantity += item.Quantity
			continue
		} else {
			itemList.TreasuryItemsQuantity += item.Quantity
		}
		if existing, ok := items[item.ItemID]; ok {
			existing.Quantity += item.Quantity
		} else {
			items[item.ItemID] = &usecase.TreasuryItem{
				ItemID:   item.ItemID,
				Name:     item.Name,
				IconURL:  item.IconURL,
				Quantity: item.Quantity,
			}
		}
	}
	itemList.Items = items
	return &itemList, nil
}
