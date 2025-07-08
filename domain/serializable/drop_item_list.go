package serializable

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DropItemList []DropItem

// DropItem make this entity serializable
type DropItem struct {
	ItemID int64  `json:"item_id" yaml:"id"`
	Rate   string `json:"rate" yaml:"rate"`
}

// Scan implementation of sql.scanner for reading db
func (d *DropItemList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid type assertion")
	}
	return json.Unmarshal(bytes, d)
}

// Value implements for writing in db
func (d DropItemList) Value() (driver.Value, error) {
	return json.Marshal(d)
}
