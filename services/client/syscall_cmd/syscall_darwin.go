package syscall_cmd

import "syscall"

func GetSyscallCmdLine(cmd string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
