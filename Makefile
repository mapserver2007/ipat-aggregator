.PHONY: gen-wire

gen-wire:
	wire gen di/wire.go

build:
	go build -o bin/ipat-aggreagtor cmd/main.go