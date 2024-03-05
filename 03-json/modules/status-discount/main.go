package main

// #include <stdlib.h>
import "C"

import (
	"encoding/json"
	"unsafe"
)

type orderSummaryInput struct {
	UserStatus int     `json:"user_status"`
	ItemCount  int     `json:"item_count"`
	Total      float64 `json:"total"`
}

type orderSummaryOutput struct {
	Code   int    `json:"code"`
	Active bool   `json:"active"`
	Reward string `json:"reward"`
}

type orderSummary struct {
	Input  orderSummaryInput    `json:"input"`
	Output []orderSummaryOutput `json:"output"`
}

func main() {}

//export calc
func calc(ptr, size uint32) (ptrSize uint64) {
	orderSummaryBytes := ptrToBytes(ptr, size)
	os := orderSummary{}
	err := json.Unmarshal(orderSummaryBytes, &os)
	if err != nil {
		panic(err)
	}

	os.Output = append(os.Output, orderSummaryOutput{
		Code:   3,
		Active: os.Input.UserStatus > 5,
		Reward: "10% off for Rouge status",
	})

	newOrderSummaryBytes, err := json.Marshal(os)
	if err != nil {
		panic(err)
	}

	ptr, size = bytesToLeakedPtr(newOrderSummaryBytes)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

// ptrToString returns a string from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToBytes(ptr uint32, size uint32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

// bytesToLeakedPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
// The pointer is not automatically managed by TinyGo hence it must be freed by the host.
func bytesToLeakedPtr(s []byte) (uint32, uint32) {
	size := C.ulong(len(s))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), s)
	return uint32(uintptr(ptr)), uint32(size)
}
