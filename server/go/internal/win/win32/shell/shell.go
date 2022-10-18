//go:build windows && amd64

package shell

import (
	"golang.org/x/sys/windows"
)

var (
	moduserenv = windows.NewLazySystemDLL("userenv.dll")
)
