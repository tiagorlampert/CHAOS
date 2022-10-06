package internal

import (
	"time"
)

const (
	TimeoutDuration   = time.Second * 30
	TimeoutExceeded   = `Timeout exceeded.`
	NoContent         = `No content.`
	TempDirectory     = `temp/`
	DatabaseDirectory = `database`
)
