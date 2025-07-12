package serializable

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// S3Screenshot serialize data about s3 storage held file data on ref {Bucket}+{Key}
type S3Screenshot struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
	S3Name string `json:"s3_name"`
}

func (s *S3Screenshot) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion .([]byte) failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s S3Screenshot) Value() (driver.Value, error) {
	return json.Marshal(s)
}
