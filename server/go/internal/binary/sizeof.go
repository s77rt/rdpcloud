//go:build windows && amd64

package binary

import "unsafe"

// SizeofUint8FromPtr takes a pointer to a uint8 sequence and returns the corresponding memory size in bytes.
// If the pointer is nil, it returns zero. It assumes that the uint8 sequence is terminated
// at a zero word; if the zero word is not present, the program may crash.
func SizeofUint8FromPtr(p *uint8) uintptr {
	if p == nil {
		return 0
	}
	var n uintptr = 0
	end := unsafe.Pointer(p)
	for {
		n++
		if *(*uint8)(end) == 0 {
			break
		}
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
	}
	return n
}

// SizeofUint16FromPtr takes a pointer to a uint16 sequence and returns the corresponding memory size in bytes.
// If the pointer is nil, it returns zero. It assumes that the uint16 sequence is terminated
// at a zero word; if the zero word is not present, the program may crash.
func SizeofUint16FromPtr(p *uint16) uintptr {
	if p == nil {
		return 0
	}
	var n uintptr = 0
	end := unsafe.Pointer(p)
	for {
		n += 2
		if *(*uint16)(end) == 0 {
			break
		}
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
	}
	return n
}
