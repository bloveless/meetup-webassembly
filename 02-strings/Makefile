.PHONY: all
all: rom-alice rom-mary rom-bob run

.PHONY: rom-alice
rom-alice:
	mkdir -p build
	rm -f ./build/alice.wasm
	cd ./modules/alice && \
	 tinygo build -o ../../build/alice.wasm -target=wasi main.go

.PHONY: rom-bob
rom-bob:
	mkdir -p build
	rm -f ./build/bob.wasm
	cd ./modules/bob && \
	 tinygo build -o ../../build/bob.wasm -target=wasi main.go

.PHONY: rom-mary
rom-mary:
	mkdir -p build
	rm -f ./build/mary.wasm
	cd ./modules/mary && \
	 tinygo build -o ../../build/mary.wasm -target=wasi main.go

.PHONY: run
run:
	go run main.go 123432456765435678
