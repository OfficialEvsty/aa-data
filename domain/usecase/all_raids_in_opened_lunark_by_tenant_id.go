package usecase

import (
	"github.com/google/uuid"
	"time"
)

// RaidDTO
// ProviderID is user who provided this raid entry
type RaidDTO struct {
	ID               uuid.UUID `json:"id"`
	ProviderID       uuid.UUID `json:"provider_id"`
	ParticipantCount uint      `json:"participant_count"`
	EventIDs         []int     `json:"event_ids"`
	Header           string    `json:"header"`
	OccurredAt       time.Time `json:"occurred_at"`
	Attendance       uint      `json:"attendance_percent"`
}

type EventDTO struct {
	TemplateID int    `json:"template_id"`
	Name       string `json:"name"`
}

type LunarkDTO struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type AllRaidsByTenantDTO struct {
	Lunark     *LunarkDTO             `json:"lunark"`
	Raids      map[uuid.UUID]*RaidDTO `json:"raids"`
	Events     map[int]*EventDTO      `json:"events"`
	Attendance uint                   `json:"attendance_percent"`
}
