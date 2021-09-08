package utils

import "time"

// GetTimeWith returns time now with the specified duration
func GetTimeWith(kind time.Duration, duration int) time.Time {
	return time.Now().Add(kind * time.Duration(duration))
}
