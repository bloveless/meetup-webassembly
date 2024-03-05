.PHONY: all
all: rom-free-item rom-free-shipping rom-status-discount run

.PHONY: rom-free-item
rom-free-item:
	mkdir -p build
	rm -f ./build/free-item.wasm
	cd ./modules/free-item && \
	 tinygo build -o ../../build/free-item.wasm -target=wasi main.go

.PHONY: rom-free-shipping
rom-free-shipping:
	mkdir -p build
	rm -f ./build/free-shipping.wasm
	cd ./modules/free-shipping && \
	 tinygo build -o ../../build/free-shipping.wasm -target=wasi main.go

.PHONY: rom-status-discount
rom-status-discount:
	mkdir -p build
	rm -f ./build/status-discount.wasm
	cd ./modules/status-discount && \
	 tinygo build -o ../../build/status-discount.wasm -target=wasi main.go

.PHONY: run
run:
	go run main.go
