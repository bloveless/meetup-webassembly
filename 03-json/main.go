package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed build/free-item.wasm
var freeItemMod []byte

//go:embed build/free-shipping.wasm
var freeShippingMod []byte

//go:embed build/status-discount.wasm
var statusDiscountMod []byte

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

func main() {
	ctx := context.Background()

	r := wazero.NewRuntime(ctx)
	defer func() {
		if err := r.Close(ctx); err != nil {
			log.Panicf("failed to close runtime: %v", err)
		}
	}()

	// Instantiate WASI, which implements host functions needed for TinyGo to
	// implement `panic`.
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	os := orderSummary{
		Input: orderSummaryInput{
			UserStatus: 1,
			ItemCount:  10,
			Total:      100.56,
		},
	}

	for _, mod := range [][]byte{freeItemMod, freeShippingMod, statusDiscountMod} {
		var err error
		os, err = runMod(ctx, r, mod, os)
		if err != nil {
			log.Panicf("unable to run order modifier: %v", err)
		}
	}

	prettyOut, err := json.MarshalIndent(os, "", "  ")
	if err != nil {
		log.Panicf("Unable to pretty print order summary")
	}

	fmt.Printf("Final order summary: %+v\n", string(prettyOut))
}

func runMod(ctx context.Context, r wazero.Runtime, modRaw []byte, os orderSummary) (orderSummary, error) {
	mod, err := r.Instantiate(ctx, modRaw)
	if err != nil {
		return orderSummary{}, fmt.Errorf("failed to instantiate module: %v", err)
	}

	// Get references to WebAssembly functions we'll use in this example.
	calc := mod.ExportedFunction("calc")
	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	osBytes, err := json.Marshal(os)
	if err != nil {
		log.Panicf("Order summary failed to marshal")
	}

	osSize := uint64(len(osBytes))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, osSize)
	if err != nil {
		log.Panicln(err)
	}

	posPtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, posPtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(uint32(posPtr), osBytes) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			posPtr, osSize, mod.Memory().Size())
	}

	// Finally, we get the greeting message "greet" printed. This shows how to
	// read-back something allocated by TinyGo.
	newOrderSummaryPtrSize, err := calc.Call(ctx, posPtr, osSize)
	if err != nil {
		log.Panicln(err)
	}

	newOrderSummaryPtr := uint32(newOrderSummaryPtrSize[0] >> 32)
	newOrderSummarySize := uint32(newOrderSummaryPtrSize[0])

	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	if newOrderSummaryPtr != 0 {
		defer func() {
			_, err := free.Call(ctx, uint64(newOrderSummaryPtr))
			if err != nil {
				log.Panicln(err)
			}
		}()
	}

	// The pointer is a linear memory offset, which is where we write the name.
	bytes, ok := mod.Memory().Read(newOrderSummaryPtr, newOrderSummarySize)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
			newOrderSummaryPtr, newOrderSummarySize, mod.Memory().Size())
	}

	newOrderSummary := orderSummary{}
	err = json.Unmarshal(bytes, &newOrderSummary)
	if err != nil {
		return orderSummary{}, err
	}

	return newOrderSummary, err
}
