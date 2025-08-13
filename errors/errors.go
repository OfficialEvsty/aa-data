package errors

import "errors"

var (
	ErrorRaidVersionMismatch        = errors.New("raid version mismatch, conflict")
	ErrorRaidPartialSavedRestricted = errors.New("raid partial saved restricted")
)
