package util

import "time"

func GetDateTime() string {
	currentTime := time.Now()
	// https://golang.org/pkg/time/#example_Time_Format
	return currentTime.Format("2006-01-02-15-04-05")
}
