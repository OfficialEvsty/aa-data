package serializable

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Participant struct {
	NicknameID uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Box        [4]Point  `json:"box"`
}

// NamedConflict similar nicknames under occurrence with box area
type NamedConflict struct {
	Similar []*Participant `json:"similar"`
	Box     [4]Point       `json:"box"`
}

func (s *NamedConflict) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

// NicknamesWithConflicts result of nicknames recognition mapped with names
type NicknamesWithConflicts struct {
	Conflicts    []NamedConflict `json:"conflicts"`
	Participants []Participant   `json:"nicknames"`
}

func (s *NicknamesWithConflicts) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}
