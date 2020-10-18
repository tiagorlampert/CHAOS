package util

import "time"

func Sleep(t time.Duration) {
	time.Sleep(t * time.Second)
}

func GetDateTime() string {
	return time.Now().Format("2006-01-02-15-04-05")
}
