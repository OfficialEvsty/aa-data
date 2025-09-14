package domain

import (
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	"github.com/google/uuid"
	"time"
)

type ExcludedParticipant struct {
	SalaryID uuid.UUID `json:"salary_id"`
	ChainID  uuid.UUID `json:"chain_id"`
	Reason   string    `json:"reason"`
}

type Salary struct {
	ID            uuid.UUID                 `json:"id"`
	Fond          int                       `json:"fond"`
	MinAttendance int                       `json:"min_attendance"`
	Tax           int                       `json:"tax"`
	CreatedAt     time.Time                 `json:"created_at"`
	Status        serializable.SalaryStatus `json:"status"`
}

type Payment struct {
	SalaryID uuid.UUID `json:"salary_id"`
	ChainID  uuid.UUID `json:"chain_id"`
	Value    int       `json:"value"`
	Reason   string    `json:"reason"`
}

type LunarkSalary struct {
	SalaryID uuid.UUID `json:"salary_id"`
	LunarkID uuid.UUID `json:"lunark_id"`
	PaidAt   time.Time `json:"paid_at"`
}
