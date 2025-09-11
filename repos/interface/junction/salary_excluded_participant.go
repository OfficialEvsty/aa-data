package junction_repos

import (
	"context"
	"github.com/google/uuid"
)

type ISalaryExcludedParticipant interface {
	AddMany(ctx context.Context, salaryID uuid.UUID, chainID uuid.UUID) error
	Remove(ctx context.Context, salaryID uuid.UUID, chainID uuid.UUID) error
}
