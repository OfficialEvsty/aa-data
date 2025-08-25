package serializable

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

// Occurrence text result inside image box occurred at specified Point's area
type Occurrence struct {
	Text string   `json:"text"`
	Box  [4]Point `json:"box"`
}

func (s *Occurrence) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s Occurrence) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// OCRData represents results of ocr working with specified image under publish_id
type OCRData struct {
	ID    uuid.UUID    `json:"id"`
	Boxes []Occurrence `json:"boxes"`
}

func (s *OCRData) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s OCRData) Value() (driver.Value, error) {
	return json.Marshal(s)
}
