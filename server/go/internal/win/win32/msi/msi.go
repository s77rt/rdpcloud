//go:build windows && amd64

package msi

import (
	"golang.org/x/sys/windows"
)

var (
	modmsi = windows.NewLazySystemDLL("msi.dll")
)
