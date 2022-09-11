//go:build windows && amd64

package encode

import (
	"golang.org/x/sys/windows"
)

func UTF16PtrToString(p *uint16) string {
	return windows.UTF16PtrToString(p)
}

func UTF16PtrFromString(s string) (*uint16, error) {
	return windows.UTF16PtrFromString(s)
}
