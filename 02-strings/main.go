package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed build/alice.wasm
var aliceMod []byte

//go:embed build/bob.wasm
var bobMod []byte

//go:embed build/mary.wasm
var maryMod []byte

func main() {
	if len(os.Args) != 2 {
		log.Panicf("You must provide exactly one numeric argument")
	}

	input := os.Args[1:2]
	mods := [][]byte{aliceMod, bobMod, maryMod}

	var modOrder []int
	for _, digit := range input[0] {
		parsedDigit, err := strconv.Atoi(string(digit))
		log.Println(parsedDigit)
		if err != nil {
			log.Panicf("Your input string must consist of only digits")
		}

		modOrder = append(modOrder, parsedDigit%len(mods))
	}

	log.Println(modOrder)

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

	for pos, modIndex := range modOrder {
		if err := runMod(ctx, r, mods[modIndex], pos); err != nil {
			log.Panicf("Error running module %d: %v", modIndex, err)
		}
	}
}

func runMod(ctx context.Context, r wazero.Runtime, modRaw []byte, position int) error {
	mod, err := r.Instantiate(ctx, modRaw)
	if err != nil {
		return fmt.Errorf("failed to instantiate module: %v", err)
	}

	// Get references to WebAssembly functions we'll use in this example.
	greeting := mod.ExportedFunction("greeting")
	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	pos := fmt.Sprintf("%d -", position)
	posSize := uint64(len(pos))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, posSize)
	if err != nil {
		log.Panicln(err)
	}

	posPtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, posPtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(uint32(posPtr), []byte(pos)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			posPtr, posSize, mod.Memory().Size())
	}

	// Finally, we get the greeting message "greet" printed. This shows how to
	// read-back something allocated by TinyGo.
	greetingPtrSize, err := greeting.Call(ctx, posPtr, posSize)
	if err != nil {
		log.Panicln(err)
	}

	greetingPtr := uint32(greetingPtrSize[0] >> 32)
	greetingSize := uint32(greetingPtrSize[0])

	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	if greetingPtr != 0 {
		defer func() {
			_, err := free.Call(ctx, uint64(greetingPtr))
			if err != nil {
				log.Panicln(err)
			}
		}()
	}

	// The pointer is a linear memory offset, which is where we write the name.
	if bytes, ok := mod.Memory().Read(greetingPtr, greetingSize); !ok {
		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
			greetingPtr, greetingSize, mod.Memory().Size())
	} else {
		fmt.Println(string(bytes))
	}

	return nil
}
