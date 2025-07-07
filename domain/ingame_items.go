package domain

import (
	"github.com/google/uuid"
	"time"
)

// AAItemTemplate in-game entity
type AAItemTemplate struct {
	ID       int64  `json:"id"` // From AACodex
	Name     string `json:"name"`
	Tier     int    `json:"tier"`
	ImageURL string `json:"image_url"`
	TierURL  string `json:"tier_url"`
}

// AAStorageItem provides info about stored items
type AAStorageItem struct {
	ID         uuid.UUID `json:"id"`
	TemplateID int64     `json:"template_id"`
	StorageID  uuid.UUID `json:"storage_id"`
	StoredAt   time.Time `json:"stored_at"`
	Quantity   int64     `json:"quantity"`
}

// AAStorage guild or user storage
type AAStorage struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	OwnerType string    `json:"owner_type"`
}
