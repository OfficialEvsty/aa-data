package junction_repos

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strconv"
	"strings"
)

type RaidItemRepository struct {
	exec db.ISqlExecutor
}

func NewRaidItemRepository(exec db.ISqlExecutor) *RaidItemRepository {
	return &RaidItemRepository{exec}
}

func (r *RaidItemRepository) AddOrUpdateItems(ctx context.Context, raidID uuid.UUID, drops []*serializable.DropItem) error {

	valueStrings := make([]string, 0, len(drops))
	valueArgs := make([]interface{}, 0, len(drops)*2)

	for i, drop := range drops {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d)", i*3+1, i*3+2, i*3+3))
		rateInt, err := strconv.Atoi(drop.Rate)
		if err != nil {
			return err
		}
		valueArgs = append(valueArgs, raidID, drop.ItemID, rateInt)
	}

	stmt := fmt.Sprintf("INSERT INTO raid_items (raid_id, item_id, rate) VALUES %s ON CONFLICT (raid_id, item_id) DO UPDATE SET rate=EXCLUDED.rate", strings.Join(valueStrings, ","))
	_, err := r.exec.ExecContext(ctx, stmt, valueArgs...)
	return err
}
func (r *RaidItemRepository) RemoveItems(ctx context.Context, raidID uuid.UUID, itemIDs []int) error {
	query := `DELETE FROM raid_items WHERE raid_id = $1 AND item_id = ANY($2)`
	_, err := r.exec.ExecContext(ctx, query, raidID, pq.Array(itemIDs))
	return err
}

func (r *RaidItemRepository) GetItems(ctx context.Context, raidID uuid.UUID) ([]*serializable.DropItem, error) {
	query := `SELECT item_id, rate FROM raid_items WHERE raid_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, raidID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dropItems []*serializable.DropItem
	dropItems = make([]*serializable.DropItem, 0)
	for rows.Next() {
		var dropItem serializable.DropItem
		err = rows.Scan(
			&dropItem.ItemID,
			&dropItem.Rate,
		)
		if err != nil {
			return nil, err
		}
		dropItems = append(dropItems, &dropItem)
	}
	return dropItems, nil
}

func (r *RaidItemRepository) RemoveItemsByRaidID(ctx context.Context, raidID uuid.UUID) error {
	query := `DELETE FROM raid_items WHERE raid_id = $1`
	_, err := r.exec.ExecContext(ctx, query, raidID)
	return err
}

func (r *RaidItemRepository) WithTx(tx *sql.Tx) junction_repos.IRaidItemRepository {
	return &RaidItemRepository{tx}
}
