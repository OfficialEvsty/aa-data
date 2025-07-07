package serializable

import (
	"encoding/json"
	"fmt"
)

type DropItemList []DropItem

// DropItem make this entity serializable
type DropItem struct {
	ItemID int64  `json:"item_id"`
	Rate   string `json:"rate"`
}

// Scan implementation of sql.scanner
func (d *DropItemList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid type assertion")
	}
	return json.Unmarshal(bytes, d)
}
