package terminal

import "syscall"

func GetHideWindowParam() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{HideWindow: true}
}
