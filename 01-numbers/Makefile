.PHONY: all
all: rom-go rom-rust rom-zig run

.PHONY: rom-go
rom-go:
	rm -f ./build/go.wasm
	cd ./modules/go && \
	 tinygo build -o ../../build/go.wasm -target=wasi main.go

# this is here as a reminder
.PHONY: rust-init
rust-init:
	rustup target add wasm32-wasi

.PHONY: rom-rust
rom-rust:
	rm -f ./build/rust.wasm
	cd ./modules/rust && \
	 cargo build --release --target wasm32-wasi
	cp ./modules/rust/target/wasm32-wasi/release/rust.wasm ./build/rust.wasm

.PHONY: rom-zig
rom-zig:
	rm -f ./build/zig.wasm
	cd ./modules/zig && \
	 zig build-lib src/main.zig -target wasm32-freestanding -dynamic -rdynamic && \
	 mv ./main.wasm ../../build/zig.wasm && \
	 rm ./main.wasm.o

.PHONY: run
run:
	go run main.go
