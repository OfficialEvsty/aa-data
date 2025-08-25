package serializable

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (s *Point) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s Point) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Conflict similar nicknames under occurrence with box area
type Conflict struct {
	Similar []uuid.UUID `json:"similar"`
	Box     [4]Point    `json:"box"`
}

type Nickname struct {
	ID  uuid.UUID `json:"id"`
	Box [4]Point  `json:"box"`
}

func (s *Nickname) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s Nickname) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Conflict) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s Conflict) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// NicknameResultWithConflicts result of nicknames recognition
type NicknameResultWithConflicts struct {
	Conflicts []Conflict `json:"conflicts"`
	Nicknames []Nickname `json:"nicknames"`
}

func (s *NicknameResultWithConflicts) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s NicknameResultWithConflicts) Value() (driver.Value, error) {
	return json.Marshal(s)
}
