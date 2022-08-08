package lib

import (
	"fmt"
	"time"
)

func ConvertTimestampToDate(timestamp int) time.Time {
	now := time.Now().Local()
	
	return time.Date(now.Year(), now.Month(), now.Day(), timestamp / 60, timestamp % 60, 0, 0, time.Local)
}

func DaysIntoYear(date time.Time) int {
	now := time.Now().Local()
	start := time.Date(now.Year(), 1, 0, 0, 0, 0, 0, time.Local)

	diff := now.Sub(start)

	return int(diff.Hours()) / 24
}

func ConvertTimestampToString(timestamp int) string {
    return fmt.Sprintf("%02d:%02d",  timestamp / 60, timestamp % 60)
}