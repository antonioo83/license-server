package utils

import (
	"time"
)

func GetTimeFromStr(dateTime string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	time, err := time.Parse(layout, dateTime)
	if err != nil {
		return time, err
	}

	return time, nil
}
