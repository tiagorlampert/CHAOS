package system

import (
	"context"
	"os/exec"
	"time"
)

func RunCmd(cmd string, timeout time.Duration) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	osType := DetectOS()
	var cmdExec *exec.Cmd
	switch osType {
	case Windows:
		cmdExec = exec.CommandContext(ctx, "cmd", "/C", cmd)
	case Linux:
		cmdExec = exec.CommandContext(ctx, "sh", "-c", cmd)
	case Darwin:
		cmdExec = exec.CommandContext(ctx, "sh", "-c", cmd)
	default:
		return []byte("os not supported")
	}

	c, err := cmdExec.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return []byte("command deadline exceeded")
		}
		return c
	}
	return c
}
