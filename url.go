package utils

import (
	"errors"
	"fmt"
	"net/url"
)

func ChangeUrlParams(urlStr string, values ...interface{}) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr, err
	}
	if len(values) == 0 || len(values)%2 != 0 {
		return urlStr, errors.New("no values to change")
	}

	q := u.Query()
	key := ""
	for k, v := range values {
		if k%2 == 0 {
			key = fmt.Sprint(v)
		} else {
			q.Set(key, fmt.Sprint(v))
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
