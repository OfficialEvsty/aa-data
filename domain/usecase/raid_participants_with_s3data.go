package usecase

import (
	"github.com/OfficialEvsty/aa-data/domain/serializable"
)

// RaidParticipantsWithS3Data if raid resolved earlier
type RaidParticipantsWithS3Data struct {
	IssuedParticipants serializable.NicknamesWithConflicts `json:"issued_participants"`
	Snapshot           serializable.S3Screenshot           `json:"snapshot"`
}
