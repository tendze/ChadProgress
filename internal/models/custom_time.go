package models

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = strings.Trim(str, `"`)

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}

func (ct CustomTime) String() string {
	return ct.Time.Format("2009-01-02 15:04:05")
}