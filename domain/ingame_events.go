package domain

import (
	"encoding/json"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"time"
)

// AABoss archeage strong enemy
type AABoss struct {
	ID            int64                     `json:"id"`
	Name          string                    `json:"name"`
	ImageGradeURL string                    `json:"img_grade_url"`
	ImageURL      string                    `json:"image_url"`
	Loot          serializable.DropItemList `json:"loot"`
}

// AAEventTemplate in-game event with specified bosses and rewards provided in necessary timing intervals
type AAEventTemplate struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Metadata json.RawMessage `json:"metadata"`
}

// EventBosses bound specified bosses to necessary in-game event
type EventBosses struct {
	BossID  int64 `json:"boss_id"`
	EventID int   `json:"event_id"`
}

// Event shows when event occurred at and what his type is
type Event struct {
	ID         uint64    `json:"id"`
	TemplateID int       `json:"template_id"`
	OccurredAt time.Time `json:"occurred_at"`
}
