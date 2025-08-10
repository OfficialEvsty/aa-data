package usecase

import (
	"encoding/json"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
)

// RaidNicknamesAndConflictsWithS3Data returns if raid unresolved
type RaidNicknamesAndConflictsWithS3Data struct {
	NicknamesWithConflicts *json.RawMessage          `json:"nicknames_with_conflicts"`
	Snapshot               serializable.S3Screenshot `json:"snapshot"`
}

type Participant struct {
	NicknameID uuid.UUID `json:"nickname_id"`
	Name       string    `json:"name"`
}

// RaidParticipantsWithS3Data if raid resolved earlier
type RaidParticipantsWithS3Data struct {
	Participants []*Participant            `json:"participants"`
	Snapshot     serializable.S3Screenshot `json:"snapshot"`
}
