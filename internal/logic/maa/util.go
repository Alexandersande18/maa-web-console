package maa

import (
	"syscall"
	"unsafe"
)

const (
	offset = 1 //char 1 byte
	null   = 0 //end of c string is 0
)

// Keep the typed pointer alive until after the native call. A uintptr alone
// does not retain the allocation for the garbage collector.
func bytePtrFromString(s string) (*byte, error) {
	return syscall.BytePtrFromString(s)
}

func CharPtr2String(r uintptr) string {
	if r == 0 {
		return ""
	}
	data := make([]byte, 0)
	// The ABI returns a native address as a machine word. Reinterpret that word
	// once, then retain a typed pointer and use pointer arithmetic (not integers).
	for p := *(**byte)(unsafe.Pointer(&r)); *p != null; p = (*byte)(unsafe.Add(unsafe.Pointer(p), offset)) {
		data = append(data, *p)
	}
	return string(data)
}

func UintPtrFromBool(b bool) uintptr {
	if b {
		return 1
	} else {
		return 0
	}
}
