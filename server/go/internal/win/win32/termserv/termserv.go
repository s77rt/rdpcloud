//go:build windows && amd64

package termserv

import (
	"golang.org/x/sys/windows"
)

var (
	modwtsapi32 = windows.NewLazySystemDLL("wtsapi32.dll")
)
