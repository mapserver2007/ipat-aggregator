GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

.PHONY: gen-wire build

gen-wire:
	wire gen di/wire.go

lint:
	goimports -w .

build:
	go mod download
	go build -o bin/ipat-aggreagtor cmd/main.go

cache-clear:
	rm -rf ./cache/colly/*

#link:
#	sudo ln -s $HOME/go/bin/go1.xx.x /usr/local/bin/go