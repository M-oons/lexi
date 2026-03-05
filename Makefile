.PHONY: help build test
.ONESHELL:

SHELL := /bin/bash

INFO_PACKAGE := github.com/m-oons/lexi/internal/info
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.1.0)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(word 1,$(shell date -u '+%Y-%m-%dT%H:%M:%SZ %Y%m%d-%H%M%S'))

help:
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build   Build the lexi binary"
	@echo "  test    Run tests"
	@echo ""

build:
	@printf "[build] Building lexi @ %s (%s) ...\n" $(VERSION) $(COMMIT)
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build \
		-ldflags " \
			-s -w \
			-X $(INFO_PACKAGE).Version=$(VERSION) \
			-X $(INFO_PACKAGE).Commit=$(COMMIT) \
			-X $(INFO_PACKAGE).Date=$(DATE) \
		" \
		-o bin/ \
		.
	@printf "[build] Built lexi @ %s (%s) -> ./bin/lexi\n" $(VERSION) $(COMMIT)

test:
	@printf "[test] Running tests...\n"
	go test \
		-v \
		-cover \
		-covermode=atomic \
		./...
	@printf "[test] Ran tests\n"
