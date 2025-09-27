package errors

import "errors"

var (
	ErrorRaidVersionMismatch        = errors.New("raid version mismatch, conflict")
	ErrorRaidPartialSavedRestricted = errors.New("raid partial saved restricted")
	ErrorSalaryVersionMismatch      = errors.New("salary version mismatch, conflict")
	ErrorItemListEmpty              = errors.New("item list is empty")
)
