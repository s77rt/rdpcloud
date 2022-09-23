//go:build windows && amd64

package netmgmt

import (
	"syscall"
)

const MAX_PREFERRED_LENGTH = 0xffffffff // (DWORD) -1 == 2^32 - 1
const TIMEQ_FOREVER = 0xffffffff        // (unsigned __LONG32) -1 == 2^32 - 1

var (
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")
)
