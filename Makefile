# Application version encoded in all the binaries.
VERSION = 0.0.0

# Base go module name.
MODULE := $(shell cat go.mod | grep -m1 module | awk '{print $$2}')

# Git revision.
REVISION := $(shell git rev-parse HEAD)

# Some sane defaults
SHELL      = /usr/bin/env bash
BINDIR      := $(CURDIR)/bin
BINNAME     ?= dogkat
PREFIX = $(HOME)/bin

# Sources for building
SOURCES := $(shell find . -type f -name *.go) go.mod go.sum

# Golang info
GOPATH := $(shell go env GOPATH)
GOBIN := $(if $(shell go env GOBIN),$(shell go env GOBIN),$(GOPATH)/bin)

# Some flags for the build to set revision and versions
FLAGS=-trimpath -ldflags "-X $(MODULE)/pkg/constants.Version=$(VERSION) -X $(MODULE)/pkg/constants.Revision=$(REVISION)"

# The targets to build
TARGET_BINARIES := dogkat-linux-amd64 dogkat-linux-arm64 dogkat-darwin-amd64 dogkat-darwin-arm64 dogkat-windows-amd64 dogkat-windows-arm64

.PHONY: all
all: $(TARGET_BINARIES)  # Make all depend on the binaries
# Create output directories (this can be simplified)
$(BINDIR):
	mkdir -p $@

# Build linux-amd64
.PHONY: dogkat-linux-amd64
dogkat-linux-amd64: $(SOURCES) | $(BINDIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(FLAGS) -o $(BINDIR)/$@ cmd/dogkat/main.go

# Build linux-arm64
.PHONY: dogkat-linux-arm64
dogkat-linux-arm64: $(SOURCES) | $(BINDIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(FLAGS) -o $(BINDIR)/$@ cmd/dogkat/main.go

# Build darwin-amd64
.PHONY: dogkat-darwin-amd64
dogkat-darwin-amd64: $(SOURCES) | $(BINDIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(FLAGS) -o $(BINDIR)/$@ cmd/dogkat/main.go

# Build darwin-arm64
.PHONY: dogkat-darwin-arm64
dogkat-darwin-arm64: $(SOURCES) | $(BINDIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(FLAGS) -o $(BINDIR)/$@ cmd/dogkat/main.go

# Build windows-amd64
.PHONY: dogkat-windows-amd64
dogkat-windows-amd64: $(SOURCES) | $(BINDIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(FLAGS) -o $(BINDIR)/$@.exe cmd/dogkat/main.go

# Build windows-arm64
.PHONY: dogkat-windows-arm64
dogkat-windows-arm64: $(SOURCES) | $(BINDIR)
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(FLAGS) -o $(BINDIR)/$@.exe cmd/dogkat/main.go

.PHONY: install
install:
	@case $$(uname -s)-$$(uname -m) in \
		Darwin-x86_64) sudo install -o root -g root -m 0755 $(BINDIR)/dogkat-darwin-amd64 $(PREFIX)/dogkat ;; \
		Linux-x86_64) sudo install -o root -g root -m 0755 $(BINDIR)/dogkat-linux-amd64 $(PREFIX)/dogkat ;; \
		Linux-aarch64) sudo install -o root -g root -m 0755 $(BINDIR)/dogkat-linux-arm64 $(PREFIX)/dogkat ;; \
		MINGW64_NT-*) sudo install -o root -g root -m 0755 $(BINDIR)/dogkat-windows-amd64.exe $(PREFIX)/dogkat.exe ;; \
		*) echo "Unsupported platform: $$(uname -s)-$$(uname -m)" ; exit 1 ;; \
	esac

# Perform license checking.
# This must pass or you will be denied by CI.
.PHONY: license
license:
	go run hack/check_license/main.go
