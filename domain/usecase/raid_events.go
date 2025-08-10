package usecase

import "github.com/google/uuid"

type RaidEventDTO struct {
	RaidID  uuid.UUID `json:"raidID"`
	EventID int       `json:"event_id"`
	Name    string    `json:"name"`
}
