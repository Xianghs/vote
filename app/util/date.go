package util

import "time"

func TimeToFormat(t time.Time) string {
	//timeFormat := "2006-01-02 15:04:05"
	timeFormat := "20060102150405"
	return t.Format(timeFormat)
}
