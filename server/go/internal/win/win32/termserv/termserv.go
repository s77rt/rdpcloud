//go:build windows && amd64

package termserv

import "syscall"

var (
	modwtsapi32 = syscall.NewLazyDLL("wtsapi32.dll")
)
