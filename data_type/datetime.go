package data_type

import (
	"database/sql/driver"
	"errors"
	"time"
	"unsafe"
)

const (
	DatetimeLayout          = "2006-01-02 15:04:05"
	DatetimeLayoutWithSlash = "2006/01/02 15:04:05"
	DateLayout              = "2006-01-02"
	DateLayoutWithSlash     = "2006/01/02"
	DateLayout1             = "01/02/2006"
)

var (
	dateFormatList = []string{
		DatetimeLayout, DatetimeLayoutWithSlash, DateLayout, DateLayoutWithSlash, DateLayout1,
	}
)

type Datetime struct {
	Time   *time.Time
	layout string
}

// NewDatetime 不指定参数则用当前时间和默认格式，t[0]:可选，格式化的时间字符串，t[1]:可选解析的格式
func NewDatetime(t ...interface{}) *Datetime {
	if len(t) > 0 {
		d := new(Datetime)
		if len(t) >= 2 {
			if _layout, ok := t[1].(string); ok {
				d.layout = _layout
			}
		}
		if d.Scan(t[0]) == nil {
			return d
		}
	}
	_now := time.Now()
	return &Datetime{
		layout: DatetimeLayout,
		Time:   &_now,
	}
}

func (d *Datetime) SetLayout(s string) *Datetime {
	d.layout = s
	return d
}
func (d *Datetime) GetLayout() string {
	if len(d.layout) == 0 {
		return DatetimeLayout
	}
	return d.layout
}

func (d Datetime) MarshalText() ([]byte, error) {
	return d.MarshalBinary()
}
func (d *Datetime) UnmarshalText(b []byte) error {
	return d.UnmarshalBinary(b)
}

func (d Datetime) MarshalBinary() ([]byte, error) {
	s := d.String()
	return *(*[]byte)(unsafe.Pointer(&s)), nil
}
func (d *Datetime) UnmarshalBinary(b []byte) error {
	d.parse(string(b))
	return nil
}

func (d Datetime) String() string {
	if d.Time == nil {
		return ""
	}
	return d.Time.Format(d.GetLayout())
}

// Format 根据指定格式格式化
func (d Datetime) Format(s string) string {
	if d.Time == nil {
		return ""
	}
	if len(s) == 0 {
		return d.String()
	}
	return d.Time.Format(s)
}

func (d Datetime) Value() (driver.Value, error) {
	if d.Time == nil {
		return nil, nil
	}
	return d.String(), nil
}

func (d *Datetime) Scan(v interface{}) error {
	switch _v := v.(type) {
	case []byte:
		return d.parse(*(*string)(unsafe.Pointer(&_v)))
	case string:
		return d.parse(_v)
	case time.Time:
		d.Time = &_v
	default:
		return errors.New("data parse error")
	}
	return nil
}

func (d *Datetime) parse(s string) (err error) {
	var t time.Time
	if len(d.layout) > 0 {
		if t, err = time.Parse(d.layout, s); err == nil {
			d.Time = &t
			return
		}
	}
	for _, _layout := range dateFormatList {
		if t, err = time.Parse(_layout, s); err == nil {
			d.Time = &t
			d.layout = _layout
			break
		}
	}
	return
}
