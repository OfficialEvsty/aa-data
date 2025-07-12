package domain

import (
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

// Journal contains RaidRecord's mapped to specified guild
type Journal struct {
	RaidID  uuid.UUID `json:"raid_id"`
	GuildID uuid.UUID `json:"guild_id"`
}

// Raid contains info about bosses and theirs loot, and guild's members participated this
type Raid struct {
	ID           uuid.UUID                 `json:"id"`
	RaidImageURL serializable.S3Screenshot `json:"raid_image_url"`
	OccurredAt   time.Time                 `json:"occurred_at"`
	Status       serializable.Status       `json:"status"`
}

// RaidEvent mapped each raid to events which passed together
type RaidEvent struct {
	RaidID  uuid.UUID `json:"raid_id"`
	EventID uuid.UUID `json:"event_id"`
}

// Attendance shows nicknames which participating event
type Attendance struct {
	RaidID     uuid.UUID `json:"raid_id"`
	NicknameID uuid.UUID `json:"nickname_id"`
}
