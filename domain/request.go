package domain

import (
	"encoding/json"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

type Request struct {
	ID         uuid.UUID                         `json:"id"`
	Type       serializable.RequestType          `json:"type"`
	Payload    json.RawMessage                   `json:"payload"`
	Status     serializable.RequestProcessStatus `json:"status"`
	Done       bool                              `json:"done"`
	EditUserID uuid.UUID                         `json:"edit_user_id"`
	CreatedAt  time.Time                         `json:"created_at"`
	SolvedAt   *time.Time                        `json:"solved_at,omitempty"`
	RollbackAt *time.Time                        `json:"rollback_at,omitempty"`
}

type TenantRequest struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	RequestID uuid.UUID `json:"request_id"`
}
