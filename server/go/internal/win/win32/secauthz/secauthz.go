//go:build windows && amd64

package secauthz

import "syscall"

var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
)
