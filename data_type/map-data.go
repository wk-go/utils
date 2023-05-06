package data_type

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type MapInt64String map[int64]string

func (d MapInt64String) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *MapInt64String) Scan(v interface{}) (err error) {
	t := make(map[int64]string)
	switch _v := v.(type) {
	case []byte:
		err = json.Unmarshal(_v, t)
	case string:
		err = json.Unmarshal([]byte(_v), t)
	default:
		return errors.New("data parse error")
	}
	if err != nil {
		return
	}
	*d = t
	return
}
