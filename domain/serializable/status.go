package serializable

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Status) Scan(src interface{}) error {
	if val, ok := src.(string); ok {
		*s = Status(strings.ToLower(val))
		return nil
	}
	return errors.New("type assertion to string failed")
}
