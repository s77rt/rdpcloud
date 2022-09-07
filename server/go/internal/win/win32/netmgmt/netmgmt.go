//go:build windows && amd64

package netmgmt

import (
	"syscall"
)

var (
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")
)
