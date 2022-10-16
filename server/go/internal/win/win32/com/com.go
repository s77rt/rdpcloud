//go:build windows && amd64

package com

import (
	"golang.org/x/sys/windows"
)

var (
	modole32 = windows.NewLazySystemDLL("ole32.dll")
)
