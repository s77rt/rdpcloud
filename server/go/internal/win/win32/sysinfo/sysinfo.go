//go:build windows && amd64

package sysinfo

import "syscall"

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
)
