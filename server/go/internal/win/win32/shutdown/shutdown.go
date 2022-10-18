//go:build windows && amd64

package shutdown

import (
	"golang.org/x/sys/windows"
)

var (
	modadvapi32 = windows.NewLazySystemDLL("advapi32.dll")
)
