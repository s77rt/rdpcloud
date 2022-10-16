//go:build windows && amd64

package sysinfo

import (
	"golang.org/x/sys/windows"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
)
