//go:build windows && amd64

package headers

const MAX_PREFERRED_LENGTH = 0xffffffff // (DWORD) -1 == 2^32 - 1
