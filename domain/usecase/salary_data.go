package usecase

import (
	"github.com/google/uuid"
	"time"
)

type SalaryContext struct {
	LunarkID      uuid.UUID `json:"lunark_id"`
	SalaryID      uuid.UUID `json:"salary_id"`
	Fond          int64     `json:"fond"`
	MinAttendance int       `json:"min_attendance"`
	Tax           int       `json:"tax"`
	CreatedAt     time.Time `json:"created_at"`
	SubmittedBy   uuid.UUID `json:"submitted_by"`
	Status        string    `json:"status"`
	Version       int       `json:"version"`
}
