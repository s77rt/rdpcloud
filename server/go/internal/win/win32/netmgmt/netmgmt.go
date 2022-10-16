//go:build windows && amd64

package netmgmt

import (
	"golang.org/x/sys/windows"
)

const MAX_PREFERRED_LENGTH = 0xFFFFFFFF
const TIMEQ_FOREVER = 0xFFFFFFFF

var (
	modnetapi32 = windows.NewLazySystemDLL("netapi32.dll")
)
