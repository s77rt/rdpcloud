//go:build windows && amd64

package secure

import "unsafe"

// ZeroMemoryUint8FromPtr assumes that the uint8 sequence is terminated
// at a zero word; if the zero word is not present, the program may crash.
func ZeroMemoryUint8FromPtr(p *uint8) {
	end := unsafe.Pointer(p)
	for *(*uint8)(end) != 0 {
		*(*uint8)(end) = 0
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
	}
}

// ZeroMemoryUint8FromPtr assumes that the uint16 sequence is terminated
// at a zero word; if the zero word is not present, the program may crash.
func ZeroMemoryUint16FromPtr(p *uint16) {
	end := unsafe.Pointer(p)
	for *(*uint16)(end) != 0 {
		*(*uint16)(end) = 0
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
	}
}
