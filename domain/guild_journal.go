package domain

import (
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

// Journal contains RaidRecord's mapped to specified guild
type Journal struct {
	LunarkID uuid.UUID `json:"lunark_id"`
	TenantID uuid.UUID `json:"tenant_id"`
}

// Lunark time interval closest to month, contains batch of raids entities to sort it
type Lunark struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Opened    bool       `json:"opened"`
}

type LunarkRaid struct {
	LunarkID uuid.UUID `json:"lunark_id"`
	RaidID   uuid.UUID `json:"raid_id"`
}

// Raid contains info about bosses and theirs loot, and guild's members participated this
type Raid struct {
	ID         uuid.UUID           `json:"id"`
	PublishID  uuid.UUID           `json:"publish_id"`
	RaidAt     *time.Time          `json:"raid_at"`
	CreatedAt  time.Time           `json:"created_at"`
	Attendance int                 `json:"attendance"`
	Status     serializable.Status `json:"status"`
	Version    int64               `json:"version"`
}

// RaidEvent mapped each raid to events which passed together
type RaidEvent struct {
	RaidID  uuid.UUID `json:"raid_id"`
	EventID int       `json:"event_id"`
}

// RaidLoot concrete item's drop by provided raid
type RaidLoot struct {
	RaidID     uuid.UUID `json:"raid_id"`
	TemplateID int64     `json:"template_id"`
	Quantity   uint64    `json:"quantity"`
}

// Attendance shows nicknames which participating event
type Attendance struct {
	RaidID     uuid.UUID `json:"raid_id"`
	NicknameID uuid.UUID `json:"nickname_id"`
}
