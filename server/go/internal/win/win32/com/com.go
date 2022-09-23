//go:build windows && amd64

package com

import "syscall"

var (
	modole32 = syscall.NewLazyDLL("ole32.dll")
)
