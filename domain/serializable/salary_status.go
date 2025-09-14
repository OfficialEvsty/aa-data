package serializable

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type SalaryStatus string

const (
	SalaryPreparing    SalaryStatus = "preparing"
	SalaryInProcessing SalaryStatus = "processing"
	SalaryPaid         SalaryStatus = "paid"
)

func (s SalaryStatus) String() string {
	return string(s)
}

func (s SalaryStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SalaryStatus) Scan(src interface{}) error {
	if val, ok := src.(string); ok {
		*s = SalaryStatus(strings.ToLower(val))
		return nil
	}
	return errors.New("type assertion to string failed")
}
