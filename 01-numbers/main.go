package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed build/go.wasm
var goMod []byte

//go:embed build/rust.wasm
var rustMod []byte

//go:embed build/zig.wasm
var zigMod []byte

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

	input := []uint64{2, 3}

	if err := runMod(ctx, r, goMod, input); err != nil {
		log.Panicf("Run Go paniced: %v", err)
	}

	if err := runMod(ctx, r, rustMod, input); err != nil {
		log.Panicf("Run rust paniced: %v", err)
	}

	if err := runMod(ctx, r, zigMod, input); err != nil {
		log.Panicf("Run zig paniced: %v", err)
	}
}

func runMod(ctx context.Context, r wazero.Runtime, modRaw []byte, input []uint64) error {
	mod, err := r.Instantiate(ctx, modRaw)
	if err != nil {
		return fmt.Errorf("failed to instantiate module: %v", err)
	}

	act := mod.ExportedFunction("act")
	if act == nil {
		return errors.New("act function didn't exist")
	}

	res, err := act.Call(ctx, input...)
	if err != nil {
		return fmt.Errorf("failed to call exported act function: %v", err)
	}

	fmt.Printf("Result: %v\n", res)

	return nil
}
