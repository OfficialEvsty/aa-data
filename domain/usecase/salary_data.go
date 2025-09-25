package usecase

import (
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

type SalaryContext struct {
	LunarkID      uuid.UUID                 `json:"lunark_id"`
	SalaryID      *uuid.UUID                `json:"salary_id"`
	Fond          int64                     `json:"fond"`
	MinAttendance int                       `json:"min_attendance"`
	Tax           int                       `json:"tax"`
	CreatedAt     time.Time                 `json:"created_at"`
	SubmittedBy   uuid.UUID                 `json:"submitted_by"`
	Status        serializable.SalaryStatus `json:"status"`
	Version       int                       `json:"version"`
}
