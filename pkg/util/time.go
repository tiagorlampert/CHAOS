package util

import "time"

func Sleep(t time.Duration) {
	time.Sleep(t * time.Second)
}
