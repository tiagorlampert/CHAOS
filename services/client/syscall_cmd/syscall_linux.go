package syscall_cmd

import "syscall"

func GetCmdSyscall(cmd string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
