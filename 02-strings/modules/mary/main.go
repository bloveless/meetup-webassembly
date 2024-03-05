package main

// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {}

//export greeting
func greeting(ptr, size uint32) (ptrSize uint64) {
	pos := ptrToString(ptr, size)
	g := fmt.Sprintf("%s Hello from Mary", pos)
	ptr, size = stringToLeakedPtr(g)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

// ptrToString returns a string from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToString(ptr uint32, size uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

// stringToLeakedPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
// The pointer is not automatically managed by TinyGo hence it must be freed by the host.
func stringToLeakedPtr(s string) (uint32, uint32) {
	size := C.ulong(len(s))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), s)
	return uint32(uintptr(ptr)), uint32(size)
}
