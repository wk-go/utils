package data_type

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func TestNewDatetime(t *testing.T) {
	data := map[string]string{
		"2002-01-06 22:22:22": DatetimeLayout,
		"2022-10-02 11:24:32": DatetimeLayout,
		"2022-10-01":          DateLayout,
	}

	for _t, layout := range data {
		datetime := NewDatetime(_t)
		_time, err := time.Parse(layout, _t)
		if err != nil {
			t.Error(err)
		}
		if datetime.Time.Unix() != _time.Unix() {
			t.Errorf("datetime[%s],_t[%s],layout[%s]", datetime, _t, layout)
		}
	}

	// 当前时间
	datetime1 := NewDatetime()
	_timeNow := time.Now()
	_layout := "2006-01-02 15:04:05"
	if datetime1.SetLayout(_layout).String() != _timeNow.Format(_layout) {
		t.Errorf("datetime[%s]!=timeNow[%s]", datetime1, _timeNow)
	}

	datetime2 := new(Datetime)
	_timeNow = time.Now()
	if datetime2.SetLayout(_layout).String() != _timeNow.Format(_layout) {
		t.Errorf("datetime[%s]!=timeNow[%s]", datetime1, _timeNow)
	}
}

func TestDatetimeMarshalJson(t *testing.T) {
	data := map[string]string{
		"2002-01-06 22:22:22": DatetimeLayout,
		"2022-10-02 11:24:32": DatetimeLayout,
		"2022-10-01":          DateLayout,
	}

	for _t := range data {
		datetime := NewDatetime(_t)
		d, _err := json.Marshal(datetime)
		if _err != nil {
			t.Error(_err)
		}

		__t := fmt.Sprintf("%#v", _t)
		if *(*string)(unsafe.Pointer(&d)) != __t {
			t.Errorf("Err:%s --- %s\n", d, __t)
		}
		t.Logf("%s --- %s\n", d, __t)
	}
}
