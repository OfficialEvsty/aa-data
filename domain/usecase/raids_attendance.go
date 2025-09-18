package usecase

import "github.com/google/uuid"

// Represens Raid IDs and Nickname ID who presented in Raid
type RaidsAttendanceByNicknameIDs = map[uuid.UUID]map[uuid.UUID]struct{}
