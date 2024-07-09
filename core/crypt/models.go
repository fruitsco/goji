package crypt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Capsule struct {
	Data       []byte
	KeyName    string
	KeyVersion int
}

var _ = driver.Valuer(&Capsule{})

// Value implements the driver.Valuer interface
func (c Capsule) Value() (driver.Value, error) {
	return json.Marshal(c)
}

var _ = sql.Scanner(&Capsule{})

// Scan implements the sql.Scanner interface
func (c *Capsule) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	return json.Unmarshal([]byte(str), c)
}
