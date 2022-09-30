//go:build windows && amd64

package netmgmt

import (
	"syscall"
)

const MAX_PREFERRED_LENGTH = 0xFFFFFFFF
const TIMEQ_FOREVER = 0xFFFFFFFF

var (
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")
)
