MKFILEPATH = $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILEDIR = $(dir $(MKFILEPATH))
DIST = $(MKFILEDIR)/dist
MKDIR = mkdir -p
GO := $(shell which go)
VERSION ?= "pre-release"
LDFLAGS := "-X 'tidalwave/cmd.Version=$(VERSION)'"
BIN = $(DIST)/tidalwave

## Create buil directories
directories:
	$(MKDIR) $(DIST)

build-darwin-amd64: directories
	GOARCH=amd64 GOOS=darwin $(GO) build -v -ldflags=$(LDFLAGS) -o $(BIN)-darwin-amd64 main.go

build-darwin-arm64: directories
	GOARCH=arm64 GOOS=darwin $(GO) build -v -ldflags=$(LDFLAGS) -o $(BIN)-darwin-arm64 main.go

build-linux-amd64: directories
	GOARCH=amd64 GOOS=linux $(GO) build -v -ldflags=$(LDFLAGS) -o $(BIN)-linux-amd64 main.go

build-linux-arm64: directories
	GOARCH=arm64 GOOS=linux $(GO) build -v -ldflags=$(LDFLAGS) -o $(BIN)-linux-arm64 main.go

build-windows-amd64: directories
	GOARCH=amd64 GOOS=windows $(GO) build -v -ldflags=$(LDFLAGS) -o $(BIN)-windows-amd64 main.go

build: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-windows-amd64