package serializable

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type RequestProcessStatus string

const (
	StatusPending  RequestProcessStatus = "pending"
	StatusRollback RequestProcessStatus = "rollback"
	StatusAccepted RequestProcessStatus = "accepted"
	StatusDeclined RequestProcessStatus = "declined"
)

func (s RequestProcessStatus) String() string {
	return string(s)
}

func (s RequestProcessStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *RequestProcessStatus) Scan(src interface{}) error {
	if val, ok := src.(string); ok {
		*s = RequestProcessStatus(strings.ToLower(val))
		return nil
	}
	return errors.New("type assertion to string failed")
}

type RequestType string

const (
	AddUserToTenantRequest RequestType = "addUserToTenant"
	AddNewNicknameRequest  RequestType = "addNewNickname"
)

func (s RequestType) String() string {
	return string(s)
}

func (s RequestType) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *RequestType) Scan(src interface{}) error {
	if val, ok := src.(string); ok {
		*s = RequestType(strings.ToLower(val))
		return nil
	}
	return errors.New("type assertion to string failed")
}
