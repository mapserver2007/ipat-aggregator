GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin


.PHONY: install
install:
	go get github.com/mark3labs/mcp-go
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: gen-wire
gen-wire:
	wire gen di/wire.go

.PHONY: lint
lint:
	goimports -w .

.PHONY: build
build:
	go mod download
	go build -o bin/mcp-server cmd/mcp/main.go

.PHONY: cache-clear
cache-clear:
	rm -rf ./cache/colly/*

#link:
#	sudo ln -s $HOME/go/bin/go1.xx.x /usr/local/bin/go