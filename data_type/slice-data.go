package data_type

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Int64Slice []int64

func NewInt64Slice(data ...int64) *Int64Slice {
	_d := Int64Slice(NewSlice(data...))
	return &_d
}

func (d Int64Slice) Len() int {
	return len(d)
}

func (d Int64Slice) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Int64Slice) Scan(v interface{}) (err error) {
	*d, err = ScanSlice[int64](v)
	return
}

func (d *Int64Slice) Append(v ...int64) {
	*d = AppendToSlice(*d, v...)
}

type StringSlice []string

func NewStringSlice(data ...string) *StringSlice {
	_d := StringSlice(NewSlice(data...))
	return &_d
}

func (d StringSlice) Len() int {
	return len(d)
}

func (d StringSlice) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *StringSlice) Scan(v interface{}) (err error) {
	*d, err = ScanSlice[string](v)
	return
}

func (d *StringSlice) Append(v ...string) {
	*d = AppendToSlice(*d, v...)
}

type SliceType interface {
	string | int64 | int
}

func NewSlice[T SliceType](s ...T) []T {
	var _data []T
	if len(s) > 0 {
		_data = s
	} else {
		_data = make([]T, 0, 10)
	}
	return _data
}
func AppendToSlice[T SliceType](target []T, s ...T) []T {
	return append(target, s...)
}

func ScanSlice[T SliceType](v interface{}) ([]T, error) {
	var err error
	t := make([]T, 0)
	switch _v := v.(type) {
	case []byte:
		err = json.Unmarshal(_v, &t)
	case string:
		err = json.Unmarshal([]byte(_v), &t)
	default:
		err = errors.New("data parse error")
	}
	if err != nil {
		return t, err
	}
	return t, err
}
