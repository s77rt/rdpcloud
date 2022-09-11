//go:build windows && amd64

package encode

import (
	"reflect"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/text/encoding/charmap"
)

// Windows1252PtrToString takes a pointer to a uint8 sequence and returns the corresponding UTF-8 encoded string.
// If the pointer is nil, it returns the empty string. It assumes that the uint8 sequence is terminated
// at a zero word; if the zero word is not present, the program may crash.
func Windows1252PtrToString(p *uint8) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint8)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(*p))
	}

	var s []uint8
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Data = uintptr(unsafe.Pointer(p))
	h.Len = n
	h.Cap = n

	c, err := charmap.Windows1252.NewDecoder().Bytes(s)
	if err != nil {
		return ""
	}

	return string(c)
}

// Windows1252PtrFromString returns the Windows-1252 encoding of the UTF-8 string
// s, with a terminating NUL added. If s contains a NUL byte at any
// location, it returns (nil, syscall.EINVAL).
func Windows1252PtrFromString(s string) (*uint8, error) {
	if strings.IndexByte(s, 0) != -1 {
		return nil, syscall.EINVAL
	}
	c, err := charmap.Windows1252.NewEncoder().String(s)
	if err != nil {
		return nil, err
	}
	a := []uint8(c + "\x00")
	return &a[0], nil
}
