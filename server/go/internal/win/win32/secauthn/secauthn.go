//go:build windows && amd64

package secauthn

import "syscall"

var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
)
