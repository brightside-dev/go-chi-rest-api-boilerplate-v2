package util

import (
	"fmt"
	"time"
)

func ParseBirthday(birthday interface{}) (time.Time, error) {
	switch v := birthday.(type) {
	case string:
		return time.Parse("2006-01-02", v)
	case []byte:
		return time.Parse("2006-01-02", string(v))
	default:
		return time.Time{}, fmt.Errorf("unexpected type for birthday: %T", v)
	}
}

func ParseDateTime(dateTime interface{}) (time.Time, error) {
	switch v := dateTime.(type) {
	case string:
		return time.Parse("2006-01-02 15:04:05", v)
	case []byte:
		return time.Parse("2006-01-02 15:04:05", string(v))
	default:
		return time.Time{}, fmt.Errorf("unexpected type for date time: %T", v)
	}
}
